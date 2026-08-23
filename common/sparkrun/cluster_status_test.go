package sparkrun

import (
	"testing"
)

func TestClusterStatusOptionsArgs(t *testing.T) {
	cluster := "c1"
	got := mustArgs(t, []string{"cluster", "status"}, ClusterStatusOptions{
		Hosts:   []string{"a", "b"},
		Cluster: &cluster,
		DryRun:  true,
	})
	want := []string{
		"cluster", "status",
		"--hosts", "a,b",
		"--cluster", "c1",
		"--dry-run",
	}
	assertEqualArgs(t, got, want)
}

func TestClusterStatusOptionsJsonMode(t *testing.T) {
	got := mustArgs(t, []string{"cluster", "status"}, ClusterStatusOptions{Hosts: []string{"a"}})
	assertEqualArgs(t, got, []string{"cluster", "status", "--hosts", "a"})
}
