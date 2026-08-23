package sparkrun

import (
	"testing"
)

func TestClusterCheckJobOptionsArgs(t *testing.T) {
	tp := 4
	port := 8000
	name := "svc"
	checkHTTP := true
	got := mustArgs(t, []string{"cluster", "check-job", "recipe"}, ClusterCheckJobOptions{
		Hosts:           Hosts([]string{"a"}),
		TensorParallel:  &tp,
		Port:            &port,
		ServedModelName: &name,
		CheckHTTPModels: &checkHTTP,
	})
	want := []string{
		"cluster", "check-job", "recipe",
		"--hosts", "a",
		"--tensor-parallel", "4",
		"--port", "8000",
		"--served-model-name", "svc",
		"--check-http-models",
	}
	assertEqualArgs(t, got, want)
}
