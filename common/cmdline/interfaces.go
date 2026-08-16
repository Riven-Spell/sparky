// Package cmdline turns a command-line option bag into a flat
// argument slice. It is generic: it knows only how to walk a struct,
// read cmd:"..." tags in declaration order, and emit --name value
// tokens. It has no knowledge of sparkrun or any specific CLI, so it
// can be ripped out and reused by any project.
package cmdline

// ArgProvider is the escape hatch for an option field whose command
// shape cannot be expressed as a single --name value pair (e.g. a
// host-list that maps to --hosts / --hosts-file / --add-host /
// --remove-host). A type implementing ArgProvider fully controls the
// tokens its field emits. The flag names inside the returned slice
// are entirely under the implementor's control.
type ArgProvider interface {
	// Args returns the argument tokens to append for this option. It
	// may return zero, one, or many tokens; zero means "nothing to
	// emit". The flag names inside the slice are up to the implementor.
	Args(tag Tag) []string
}
