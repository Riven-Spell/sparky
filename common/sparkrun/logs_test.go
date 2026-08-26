package sparkrun

import (
	"testing"

	"github.com/Riven-Spell/sparky/common/cmdline"
)

// mustLogsArgs builds the Logs argument list from prefix + bag + "--follow",
// mirroring how [cliClient.Logs] constructs its command.
func mustLogsArgs(t *testing.T, prefix []string, bag LogsOptions) []string {
	t.Helper()
	got, err := cmdline.BuildArgs(prefix, bag, "--follow")
	if err != nil {
		t.Fatalf("BuildArgs(%T): %v", bag, err)
	}
	return got
}

func TestLogsOptionsArgs(t *testing.T) {
	tail := 200
	cluster := "mylab"
	got := mustLogsArgs(t, []string{"logs", "recipe"}, LogsOptions{
		Hosts:   Hosts([]string{"a", "b"}),
		Cluster: &cluster,
		Tail:    &tail,
	})
	assertEqualArgs(t, got, []string{
		"logs", "recipe",
		"--hosts", "a,b",
		"--cluster", "mylab",
		"--tail", "200",
		"--follow",
	})
}

func TestLogsOptionsHostsFromFile(t *testing.T) {
	got := mustLogsArgs(t, []string{"logs", "recipe"}, LogsOptions{
		Hosts: HostsFromFile("/tmp/hosts"),
	})
	assertEqualArgs(t, got, []string{"logs", "recipe", "--hosts-file", "/tmp/hosts", "--follow"})
}

func TestLogsOptionsAllSources(t *testing.T) {
	got := mustLogsArgs(t, []string{"logs", "sparkrun_abc"}, LogsOptions{
		AllSources: true,
	})
	assertEqualArgs(t, got, []string{"logs", "sparkrun_abc", "--all-sources", "--follow"})
}

func TestLogsOptionsDefaults(t *testing.T) {
	got := mustLogsArgs(t, []string{"logs", "recipe"}, LogsOptions{})
	assertEqualArgs(t, got, []string{"logs", "recipe", "--follow"})
}
