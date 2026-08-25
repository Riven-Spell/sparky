package sparkrun

import (
	"sort"

	"github.com/Riven-Spell/sparky/common/cmdline"
)

// KeyValueOpts is a repeatable set of --<name> key=value options for a
// single command (e.g. --executor-opt privileged=false). Use it for any
// repeatable key=value option: cluster --executor-opt, run -o/-e, and
// so on. Keys are emitted in sorted order so argument output is
// deterministic regardless of map iteration order.
type KeyValueOpts map[string]string

// Args emits one --<name> key=value pair per entry, using the flag name
// from the field's cmd:"..." tag. An empty value emits --<name> key=.
func (o KeyValueOpts) Args(tag cmdline.Tag) []string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	args := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		args = append(args, "--"+tag.Name, k+"="+o[k])
	}
	return args
}

var _ cmdline.ArgProvider = KeyValueOpts(nil)
