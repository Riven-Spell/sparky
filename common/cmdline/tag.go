package cmdline

import (
	"fmt"
	"strings"
)

// Tag is the parsed form of a `cmd:"..."` struct tag.
//
//   - Name is the flag name (the segment before the first comma).
//
// The remaining comma-separated segments are options:
//
//   - a bare identifier is a custom method name (see [Tag.Method]).
//   - `default=X` sets the value used when the field is its zero value.
type Tag struct {
	Name    string
	Method  string
	Default string
}

// ParseTag parses a `cmd` struct tag and returns its [Tag]. The tag
// must be non-empty and start with a flag name. A malformed tag (no
// flag name, or more than one custom method) is returned as an error.
func ParseTag(tag string) (Tag, error) {
	flagName, options := tag, ""
	if before, after, ok := strings.Cut(tag, ","); ok {
		flagName, options = before, after
	}
	if flagName == "" {
		return Tag{}, fmt.Errorf("cmdline: empty flag name in cmd tag %q", tag)
	}
	parsed := Tag{Name: flagName}
	if options == "" {
		return parsed, nil
	}
	for _, option := range strings.Split(options, ",") {
		switch {
		case option == "":
			continue
		case strings.HasPrefix(option, "default="):
			parsed.Default = option[len("default="):]
		default:
			if parsed.Method != "" {
				return Tag{}, fmt.Errorf("cmdline: multiple custom methods in cmd tag %q", tag)
			}
			parsed.Method = option
		}
	}
	return parsed, nil
}
