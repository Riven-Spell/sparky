package sparkrun

import (
	"testing"

	"github.com/Riven-Spell/sparky/common/cmdline"
)

func TestClusterUpdateOptionsArgs(t *testing.T) {
	ring := ETopology.Ring()
	got := mustArgs(t, nil, ClusterUpdateOptions{
		Hosts:       RemoveHosts([]string{"a", "b"}),
		Description: new("d"),
		Topology:    &ring,
	})
	want := []string{
		"--remove-host", "a,b",
		"--description", "d",
		"--topology", "ring",
	}
	assertEqualArgs(t, got, want)
}

func TestClusterUpdateOptionsHostDiffModes(t *testing.T) {
	got := mustArgs(t, nil, ClusterUpdateOptions{Hosts: Hosts([]string{"x", "y", "z"})})
	assertEqualArgs(t, got, []string{"--hosts", "x,y,z"})

	got = mustArgs(t, nil, ClusterUpdateOptions{Hosts: HostsFromFile("/f")})
	assertEqualArgs(t, got, []string{"--hosts-file", "/f"})

	got = mustArgs(t, nil, ClusterUpdateOptions{Hosts: AddHosts([]string{"x", "y", "z"})})
	assertEqualArgs(t, got, []string{"--add-host", "x,y,z"})

	got = mustArgs(t, nil, ClusterUpdateOptions{Hosts: RemoveHosts([]string{"x", "y", "z"})})
	assertEqualArgs(t, got, []string{"--remove-host", "x,y,z"})

	got = mustArgs(t, nil, ClusterUpdateOptions{Hosts: CompoundHostDiff([]string{"x", "y"}, []string{"a", "b"})})
	assertEqualArgs(t, got, []string{"--add-host", "x,y", "--remove-host", "a,b"})
}

func TestClusterUpdateOptionsNilHosts(t *testing.T) {
	got := mustArgs(t, nil, ClusterUpdateOptions{})
	assertEqualArgs(t, got, []string{})
}

func mustArgs(t *testing.T, prefix []string, bag any) []string {
	t.Helper()
	got, err := cmdline.BuildArgs(prefix, bag)
	if err != nil {
		t.Fatalf("BuildArgs(%T): %v", bag, err)
	}
	return got
}

func assertEqualArgs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len mismatch\ngot:  %#v\nwant: %#v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d diff: got %q want %q\ngot:  %#v\nwant: %#v", i, got[i], want[i], got, want)
		}
	}
}
