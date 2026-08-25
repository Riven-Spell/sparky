package sparkrun

import (
	"testing"
)

func TestClusterCreateOptionsArgs(t *testing.T) {
	desc := "hello"
	mode := ETransferMode.Delegated()
	inter := ETransferInterface.Mgmt()
	exec := "docker"
	sched := "greedy"
	mem := 0.85
	got := mustArgs(t, nil, ClusterCreateOptions{
		Hosts:             Hosts([]string{"a", "b"}),
		Description:       &desc,
		User:              strPtr("u"),
		CacheDir:          strPtr("/cache"),
		TransferMode:      &mode,
		TransferInterface: &inter,
		SetDefault:        true,
		Executor:          &exec,
		ExecutorOpts:      KeyValueOpts{"privileged": "false", "shm_size": "16g"},
		Scheduler:         &sched,
		MaxGPUMemUtil:     &mem,
	})
	want := []string{
		"--hosts", "a,b",
		"--description", "hello",
		"--user", "u",
		"--cache-dir", "/cache",
		"--transfer-mode", "delegated",
		"--transfer-interface", "mgmt",
		"--default",
		"--executor", "docker",
		"--executor-opt", "privileged=false",
		"--executor-opt", "shm_size=16g",
		"--scheduler", "greedy",
		"--max-gpu-mem-util", "0.85",
	}
	assertEqualArgs(t, got, want)
}

func TestClusterCreateOptionsMinimal(t *testing.T) {
	got := mustArgs(t, nil, ClusterCreateOptions{Hosts: Hosts([]string{"a"})})
	want := []string{"--hosts", "a"}
	assertEqualArgs(t, got, want)
}

func strPtr(s string) *string { return &s }
