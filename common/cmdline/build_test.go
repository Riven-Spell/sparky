package cmdline

import (
	"strconv"
	"testing"
)

func TestBuildArgsKinds(t *testing.T) {
	type bag struct {
		Val   int      `cmd:"val"`
		Flag  bool     `cmd:"flag"`
		Tests []string `cmd:"tests"`
		Ptr   *int     `cmd:"ptr"`
	}

	ptrVal := 7
	got, err := BuildArgs([]string{"verb", "sub"}, bag{
		Val:   5,
		Flag:  true,
		Tests: []string{"a", "b"},
		Ptr:   &ptrVal,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{
		"verb", "sub",
		"--val", "5",
		"--flag",
		"--tests", "a,b",
		"--ptr", "7",
	}
	assertArgs(t, got, want)
}

func TestBuildArgsSkipsEmpty(t *testing.T) {
	type bag struct {
		Str   string   `cmd:"str"`
		Bool  bool     `cmd:"bool"`
		Slice []string `cmd:"slice"`
		Ptr   *int     `cmd:"ptr"`
	}
	got, err := BuildArgs([]string{"verb"}, bag{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"verb"}
	assertArgs(t, got, want)
}

func TestBuildArgsDeclarationOrder(t *testing.T) {
	type bag struct {
		Z string `cmd:"z"`
		A string `cmd:"a"`
		M string `cmd:"m"`
	}
	got, err := BuildArgs(nil, bag{Z: "1", A: "2", M: "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"--z", "1", "--a", "2", "--m", "3"}
	assertArgs(t, got, want)
}

type providerBag struct {
	Val string   `cmd:"val"`
	Cmd provider `cmd:"cmd"`
}

type provider struct {
	multi bool
}

func (p provider) Args(Tag) []string {
	if !p.multi {
		return nil
	}
	return []string{"--p", "1", "--p", "2"}
}

func TestBuildArgsArgProvider(t *testing.T) {
	got, err := BuildArgs(nil, providerBag{Val: "v", Cmd: provider{multi: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"--val", "v", "--p", "1", "--p", "2"}
	assertArgs(t, got, want)

	got, err = BuildArgs(nil, providerBag{Val: "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want = []string{"--val", "v"}
	assertArgs(t, got, want)
}

type customBag struct {
	Base int `cmd:"base"`
	X    int `cmd:"x,Val"`
}

func (b *customBag) Val() []string {
	return []string{"--x", "double:" + strconv.Itoa(b.X*2)}
}

func TestBuildArgsCustomMethod(t *testing.T) {
	got, err := BuildArgs(nil, customBag{Base: 1, X: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"--base", "1", "--x", "double:6"}
	assertArgs(t, got, want)
}

type missingBag struct {
	X int `cmd:"x,NoSuchMethod"`
}

func TestBuildArgsCustomMethodMissingError(t *testing.T) {
	_, err := BuildArgs(nil, missingBag{X: 1})
	if err == nil {
		t.Fatal("expected error for missing custom method")
	}
}

func TestBuildArgsNonStructErrors(t *testing.T) {
	for _, bag := range []any{
		42,
		"str",
		[]string{"a"},
		&[]string{"a"},
		map[string]int{"k": 1},
	} {
		if _, err := BuildArgs(nil, bag); err == nil {
			t.Errorf("expected error for non-struct bag %T", bag)
		}
	}
}

func TestBuildArgsNilPtrOK(t *testing.T) {
	type bag struct {
		Val   int      `cmd:"val"`
		Flag  bool     `cmd:"flag"`
		Tests []string `cmd:"tests"`
		Ptr   *int     `cmd:"ptr"`
	}
	var ptr *bag
	got, err := BuildArgs([]string{"verb"}, ptr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"verb"}
	assertArgs(t, got, want)
}

type defBag struct {
	Port *int     `cmd:"port,default=8080"`
	Host *string  `cmd:"host,default=localhost"`
	Tags []string `cmd:"tags,default=foo"`
}

func TestBuildArgsDefault(t *testing.T) {
	got, err := BuildArgs(nil, defBag{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"--port", "8080", "--host", "localhost", "--tags", "foo"}
	assertArgs(t, got, want)

	got, err = BuildArgs(nil, defBag{Port: new(9090), Host: new("example.com")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want = []string{"--port", "9090", "--host", "example.com", "--tags", "foo"}
	assertArgs(t, got, want)

	got, err = BuildArgs(nil, defBag{Tags: []string{"a", "b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want = []string{"--port", "8080", "--host", "localhost", "--tags", "a,b"}
	assertArgs(t, got, want)
}

type defScalarBag struct {
	Host string `cmd:"host,default=localhost"`
}

func TestBuildArgsDefaultNonNullable(t *testing.T) {
	if _, err := BuildArgs(nil, defScalarBag{}); err == nil {
		t.Fatal("expected error for default on non-nullable field")
	}
}

func TestParseTag(t *testing.T) {
	testCases := []struct {
		input string
		want  Tag
	}{
		{"port", Tag{Name: "port"}},
		{"x,Val", Tag{Name: "x", Method: "Val"}},
		{"port,default=8080", Tag{Name: "port", Default: "8080"}},
		{"port,Val,default=8080", Tag{Name: "port", Method: "Val", Default: "8080"}},
	}
	for _, testCase := range testCases {
		got, err := ParseTag(testCase.input)
		if err != nil {
			t.Errorf("ParseTag(%q): unexpected error %v", testCase.input, err)
			continue
		}
		if got != testCase.want {
			t.Errorf("ParseTag(%q)=%+v want %+v", testCase.input, got, testCase.want)
		}
	}
}

func TestParseTagErrors(t *testing.T) {
	for _, input := range []string{",x", "a,Val,Other"} {
		if _, err := ParseTag(input); err == nil {
			t.Errorf("ParseTag(%q): expected error", input)
		}
	}
}

func assertArgs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len mismatch\ngot:  %#v\nwant: %#v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q want %q\ngot:  %#v\nwant: %#v", i, got[i], want[i], got, want)
		}
	}
}
