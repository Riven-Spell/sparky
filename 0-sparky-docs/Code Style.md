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
Golang does not have native support for enums. Thus, an external package must be used to mimic them. My package of choice would be my own enum package, `github.com/Riven-Spell/enum`. Enums are best treated as their own package in the project, as they are not likely to reference other packages and are easily a primitive.
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