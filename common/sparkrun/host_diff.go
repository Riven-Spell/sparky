package sparkrun

import (
	"strings"

	"github.com/Riven-Spell/sparky/common/cmdline"
)

// HostsList describes a host-list replacement: either a fixed list of
// addresses ([Hosts]) or a file of addresses ([HostsFromFile]). Verbs
// that only ever replace the whole host set (e.g. `cluster create`)
// accept a HostsList and nothing else.
//
// HostsList is a [cmdline.ArgProvider]: each implementation fully
// controls the tokens its option field emits. The unexported hostsList()
// method seals the interface, so only types defined inside this package
// can satisfy it.
type HostsList interface {
	HostsDiff
	hostsList()
}

// HostsDiff describes any host-list mutation accepted by verbs like
// `sparkrun cluster update`: a replacement ([HostsList], i.e. [Hosts] /
// [HostsFromFile]) or an incremental add/remove ([AddHosts],
// [RemoveHosts], [CompoundHostDiff]). The caller is responsible for
// passing at most one HostsDiff per call (sparkrun itself rejects
// combinations of a replacement with an incremental edit).
//
// The unexported hostDiff() method seals the interface, so any type
// satisfying HostsDiff must be defined inside this package; an arbitrary
// [cmdline.ArgProvider] from outside is not accepted as a HostsDiff.
type HostsDiff interface {
	cmdline.ArgProvider
	hostDiff()
}

type hostsReplace struct {
	values []string
}

var _ HostsList = hostsReplace{}
var _ HostsDiff = hostsReplace{}

func (h hostsReplace) Args(cmdline.Tag) []string {
	if len(h.values) == 0 {
		return nil
	}
	joined := strings.Join(h.values, ",")
	return []string{"--hosts", joined}
}

func (hostsReplace) hostsList() {}
func (hostsReplace) hostDiff()  {}

// Hosts replaces the cluster's host list with the given addresses.
// Emits `--hosts a,b,c`.
func Hosts(values []string) HostsList {
	return hostsReplace{values: values}
}

type hostsFromFile struct {
	path string
}

var _ HostsList = hostsFromFile{}
var _ HostsDiff = hostsFromFile{}

func (h hostsFromFile) Args(cmdline.Tag) []string {
	if h.path == "" {
		return nil
	}
	return []string{"--hosts-file", h.path}
}

func (hostsFromFile) hostsList() {}
func (hostsFromFile) hostDiff()  {}

// HostsFromFile replaces the cluster's host list with one address per
// line of path. Emits `--hosts-file <path>`.
func HostsFromFile(path string) HostsList {
	return hostsFromFile{path: path}
}

type hostsAdd struct {
	values []string
}

var _ HostsDiff = hostsAdd{}

func (h hostsAdd) Args(cmdline.Tag) []string {
	if len(h.values) == 0 {
		return nil
	}
	return []string{"--add-host", strings.Join(h.values, ",")}
}

func (hostsAdd) hostDiff() {}

// AddHosts appends the given addresses to the cluster. Emits
// `--add-host a,b,c` (sparkrun accepts a single repeat with
// comma-separated values).
func AddHosts(values []string) HostsDiff {
	return hostsAdd{values: values}
}

type hostsRemove struct {
	values []string
}

var _ HostsDiff = hostsRemove{}

func (h hostsRemove) Args(cmdline.Tag) []string {
	if len(h.values) == 0 {
		return nil
	}
	return []string{"--remove-host", strings.Join(h.values, ",")}
}

func (hostsRemove) hostDiff() {}

// RemoveHosts removes the given addresses from the cluster. Emits
// `--remove-host a,b,c` (sparkrun accepts a single repeat with
// comma-separated values).
func RemoveHosts(values []string) HostsDiff {
	return hostsRemove{values: values}
}

type compoundHostDiff struct {
	added   hostsAdd
	removed hostsRemove
}

var _ HostsDiff = compoundHostDiff{}

func (c compoundHostDiff) Args(tag cmdline.Tag) []string {
	args := c.added.Args(tag)
	args = append(args, c.removed.Args(tag)...)
	return args
}

func (compoundHostDiff) hostDiff() {}

// CompoundHostDiff performs an incremental edit in one call, mirroring
// how sparkrun's `cluster update` accepts --add-host and --remove-host
// together. Emits the add then remove flags (each a single
// comma-separated repeat).
func CompoundHostDiff(added, removed []string) HostsDiff {
	return compoundHostDiff{
		added:   hostsAdd{values: added},
		removed: hostsRemove{values: removed},
	}
}
