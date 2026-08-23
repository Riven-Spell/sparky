package sparkrun

import (
	"testing"
)

func TestClusterCreateOptionsArgs(t *testing.T) {
	desc := "hello"
	mode := ETransferMode.Delegated()
	inter := ETransferInterface.Mgmt()
	got := mustArgs(t, nil, ClusterCreateOptions{
		Hosts:             Hosts([]string{"a", "b"}),
		Description:       &desc,
		User:              strPtr("u"),
		CacheDir:          strPtr("/cache"),
		TransferMode:      &mode,
		TransferInterface: &inter,
		SetDefault:        true,
	})
	want := []string{
		"--hosts", "a,b",
		"--description", "hello",
		"--user", "u",
		"--cache-dir", "/cache",
		"--transfer-mode", "delegated",
		"--transfer-interface", "mgmt",
		"--default",
	}
	assertEqualArgs(t, got, want)
}

func TestClusterCreateOptionsMinimal(t *testing.T) {
	got := mustArgs(t, nil, ClusterCreateOptions{Hosts: Hosts([]string{"a"})})
	want := []string{"--hosts", "a"}
	assertEqualArgs(t, got, want)
}

func strPtr(s string) *string { return &s }
