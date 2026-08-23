This document mainly exists to inform AI agents about rules to follow when writing code. However, it may prove relevant to human contributors as well.
# One file, one purpose.
Big code files are hard to decipher what purpose they serve, to both human operators and language models. One file should serve one purpose.

For instance, if you have 4 commands, `<root>`, `foo`, `foo bar`, and `baz`, they shouldn't all be in the same `commands.go` file. They should be in `cmd/root.go`, `cmd/foo.go`, `cmd/foo_bar.go`, and `cmd/baz.go` respectively.
# Import loops
Golang cannot handle recursive import loops. Thus, the use of a `common` package is well, common. This creates a space for shared models/datastructures/etc. to exist between layers and be referenced.

If an import loop is created, simply branch the package out to produce a common package to escape the import loop.
# One package, one purpose.
Why should credential management live in the same package as primitive data structures? Why should they be in the same package as models used between services? Why would project-specific models be in the same package as models for an external dependency?

A separate package breaks it down to a logical module. `cmd` provides command line. `common`, and it's child packages provide common data structures. `agent` provides the daemon that lives on each node. `manager` provides the backend for the manager. `webui` provides for hosting the webui.
# Enums
Golang does not have native support for enums. Thus, an external package must be used to mimic them. My package of choice would be my own enum package, `https://pkg.go.dev/github.com/Riven-Spell/enum/v2`. Enums are best treated as their own package in the project, as they are not likely to reference other packages and are easily a primitive.
# Options bags
A function will rarely have more than two or three mandatory inputs. `FooBar(x, y, z, true, false, "", false, false, 50)` is insane to read and parse through. Options bags provide a named way to handle both optionality, and if it gets long enough, mandatory values.

i.e.
```go
type FooBarOptions struct {
	Mandatory int // no pointer means mandatory!
	
	// Create a partition between mandatory and optional-- this makes the human brain observe a difference.
	
	Baz *string // pointer means optional!
	Boop []string // but be wary, lists,
	Befuddle map[string]string // and maps are inherently nullable, and do not need a pointer. Making these pointers requires a double nil check-- this is unnatural.
}
```

If there are no mandatory parameters inside the options bag, make the options bag a vararg and accept the first value present. Make no effort to compress. Multiple options bags is operator error. No options bag defaults all optional values. 

Foo goes from `FooBar(<insane list>)` to `FooBar(mandatoryX, opts ...FooBarOptions)`, or if there's many mandatory parameters, just `FooBar(opts FooBarOptions)`. Many options becomes readable, because the inputs become labeled.
# Struct field labeling
Always fully label struct fields. `FooBar{x, y, z}` sure is convenient, but you can't tell me what the fields mean when it gets large enough. Same issue that spawned options bags.
# Unexported internals
Types, functions, and variables that are not referenced outside their package should remain unexported. Export only what external consumers need. This reduces the public surface area, makes refactoring easier, and avoids committing to a stable API for implementation details.
# Explicit error discarding
When an error return is genuinely not actionable (e.g. closing an HTTP response body, encoding to a response writer after the connection is dead), discard it explicitly with `_ = fn()`. Do not leave the return value dangling — this makes the intent clear and keeps linters happy.

`.Close()` is the canonical example: pipe and file descriptor closes during cleanup almost never have a meaningful recovery path, so always write `_ = f.Close()` rather than bare `f.Close()`.
# Readable variable names
Do not shorthand variable names. Every variable, parameter, and field should be fully spelled out so its meaning is legible at a glance. `ft` is not readable; `field` is.

Prefer, e.g. `bagValue` over `rv`, `addressable` over `av`, `fieldIndex` over `i`, `options` over `rest`. Loop counters and test helpers are not exempt — avoid packing meaning into single letters or two-letter abbreviations.
# Option Overloading
Golang does not natively support overloads-- i.e. Foo(x int); Foo(x func()). To work around this, utilize interfaces that allow us to dynamically satisfy requirements-- see sparkrun package for examples.
# No pointless standard-library helpers
Do not write helper functions that merely reimplement something the standard library already provides. A hand-rolled `joinComma(s []string) string` that just does `strings.Join(s, ",")` adds abstraction (and a dedicated import) with zero benefit. Use the standard library directly, and reserve helper functions for logic the library genuinely lacks.
