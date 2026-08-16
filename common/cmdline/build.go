package cmdline

import (
	"fmt"
	"reflect"
	"strings"
)

// argProviderType is the reflect type of [ArgProvider], used to detect
// fields whose type implements the interface.
var argProviderType = reflect.TypeFor[ArgProvider]()

// BuildArgs renders bag into the final argument list. prefix is
// prepended verbatim; every field of bag carrying a cmd:"..." tag is
// then appended in declaration order. It returns an error only when a
// field cannot be rendered (e.g. a malformed tag, or a custom method
// named by the tag that is missing or has the wrong signature) -- see
// [ArgProvider], [ParseTag], and the custom-method form.
func BuildArgs(prefix []string, bag any) ([]string, error) {
	args := append([]string(nil), prefix...)

	bagValue := reflect.ValueOf(bag)
	if !bagValue.IsValid() {
		return args, nil
	}
	switch bagValue.Kind() {
	case reflect.Pointer:
		if bagValue.IsNil() {
			return args, nil
		}
		bagValue = bagValue.Elem()
		if bagValue.Kind() != reflect.Struct {
			return nil, fmt.Errorf("cmdline: bag must be a struct or *struct, got %s", bagValue.Kind())
		}
	case reflect.Struct:
		// A struct passed by value is not addressable, which hides any
		// pointer-receiver methods. Re-wrap in *T so both value and
		// pointer receiver methods are reachable.
		addressable := reflect.New(bagValue.Type()).Elem()
		addressable.Set(bagValue)
		bagValue = addressable
	default:
		return nil, fmt.Errorf("cmdline: bag must be a struct or *struct, got %s", bagValue.Kind())
	}

	bagType := bagValue.Type()
	for fieldIndex := 0; fieldIndex < bagValue.NumField(); fieldIndex++ {
		field := bagType.Field(fieldIndex)
		rawTag := field.Tag.Get("cmd")
		if rawTag == "" {
			continue
		}
		parsed, err := ParseTag(rawTag)
		if err != nil {
			return nil, fmt.Errorf("cmdline: field %s: %w", field.Name, err)
		}
		tokens, err := emit(parsed, bagValue, bagValue.Field(fieldIndex))
		if err != nil {
			return nil, fmt.Errorf("cmdline: field %s: %w", parsed.Name, err)
		}
		if len(tokens) > 0 {
			args = append(args, tokens...)
		}
	}
	return args, nil
}

func emit(tag Tag, bag, fieldValue reflect.Value) ([]string, error) {
	// Escape hatch for objects providing their own argument parsing
	if fieldValue.IsValid() && fieldValue.CanInterface() && fieldValue.Type().Implements(argProviderType) {
		if (fieldValue.Kind() == reflect.Interface || fieldValue.Kind() == reflect.Pointer) && fieldValue.IsNil() {
			return nil, nil
		}
		return fieldValue.Interface().(ArgProvider).Args(tag), nil
	}

	// Escape hatch for custom emission methods on parent struct
	if tag.Method != "" {
		return callCustom(bag, tag.Method)
	}

	// A default fills in the --name value when a nullable field
	// (pointer, slice, or map) is unset. For any other field kind the
	// default carries no meaning, so it is a programmer error to set it.
	switch fieldValue.Kind() {
	case reflect.Pointer:
		if fieldValue.IsNil() {
			return withDefault(tag)
		}
		// An explicit pointer is respected as-is, never the default.
		return emitValue(tag, fieldValue.Elem())
	case reflect.Slice:
		if fieldValue.IsNil() {
			return withDefault(tag)
		}
		return emitValue(tag, fieldValue)
	case reflect.Map:
		if fieldValue.IsNil() {
			return withDefault(tag)
		}
		return emitValue(tag, fieldValue)
	default:
		if tag.Default != "" {
			return nil, fmt.Errorf("cmdline: field %s: default '%s' only allowed on nullable fields (pointer, slice, map)", tag.Name, tag.Default)
		}
		return emitValue(tag, fieldValue)
	}
}

// withDefault emits --name <default> when the tag carries one, or
// nothing when it does not. Used for unset nullable fields.
func withDefault(tag Tag) ([]string, error) {
	if tag.Default == "" {
		return nil, nil
	}
	return []string{"--" + tag.Name, tag.Default}, nil
}

// emitValue renders a present field -- bool, string, slice, map, or
// scalar -- into --name value tokens. It never consults a default.
func emitValue(tag Tag, fieldValue reflect.Value) ([]string, error) {
	switch fieldValue.Kind() {
	case reflect.Bool:
		if fieldValue.Bool() {
			return []string{"--" + tag.Name}, nil
		}
		return nil, nil
	case reflect.String:
		if fieldValue.String() == "" {
			return nil, nil
		}
		return []string{"--" + tag.Name, fieldValue.String()}, nil
	case reflect.Slice:
		if fieldValue.Len() == 0 {
			return nil, nil
		}
		return []string{"--" + tag.Name, joinSlice(fieldValue)}, nil
	case reflect.Map:
		if fieldValue.Len() == 0 {
			return nil, nil
		}
		return []string{"--" + tag.Name, joinMap(fieldValue)}, nil
	default:
		if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			return nil, nil
		}
		return []string{"--" + tag.Name, fmt.Sprint(fieldValue.Interface())}, nil
	}
}

func joinSlice(fieldValue reflect.Value) string {
	parts := make([]string, fieldValue.Len())
	for partIndex := 0; partIndex < fieldValue.Len(); partIndex++ {
		parts[partIndex] = fmt.Sprint(fieldValue.Index(partIndex))
	}
	return strings.Join(parts, ",")
}

func joinMap(fieldValue reflect.Value) string {
	parts := make([]string, 0, fieldValue.Len())
	for _, key := range fieldValue.MapKeys() {
		parts = append(parts, fmt.Sprint(key.Interface())+"="+fmt.Sprint(fieldValue.MapIndex(key).Interface()))
	}
	return strings.Join(parts, ",")
}

// callCustom looks up method on the bag being encoded (canonical
// signature `func (b *Bag) Name() []string`) and returns its tokens.
// The method reads the bag's own fields, so it is fully type-safe and
// needs no reflection in the method body. The method must be exported
// (reflect only sees exported methods). A tag naming a method that is
// missing or returns the wrong type is returned as an error.
func callCustom(bag reflect.Value, method string) ([]string, error) {
	methodValue := bag.Addr().MethodByName(method)
	if !methodValue.IsValid() {
		methodValue = bag.MethodByName(method)
	}
	if !methodValue.IsValid() {
		return nil, fmt.Errorf("no method %s() []string on option bag %s", method, bag.Type())
	}
	results := methodValue.Call(nil)
	if len(results) != 1 || !results[0].IsValid() {
		return nil, fmt.Errorf("method %s() on option bag %s must return []string", method, bag.Type())
	}
	tokens, ok := results[0].Interface().([]string)
	if !ok {
		return nil, fmt.Errorf("method %s() on option bag %s must return []string", method, bag.Type())
	}
	return tokens, nil
}
