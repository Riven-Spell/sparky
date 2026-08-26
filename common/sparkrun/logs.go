package sparkrun

import (
	"context"
	"io"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
)

// LogsOptions configures [Client.Logs].
type LogsOptions struct {
	// Hosts is the host list to query (--hosts or --hosts-file), expressed
	// as a [HostsList] ([Hosts] or [HostsFromFile]).
	Hosts HostsList `cmd:"hosts"`
	// Cluster uses a saved cluster by name (--cluster) to select the
	// host set. Mutually exclusive with Hosts.
	Cluster *string `cmd:"cluster"`
	// TensorParallel is the tensor-parallelism override (--tensor-parallel),
	// so the log source matches a run with the same override.
	TensorParallel *int `cmd:"tensor-parallel"`
	// Port is the inference-server port override (--port) used to match
	// a run-time override.
	Port *int `cmd:"port"`
	// ServedModelName is the served-model-name override
	// (--served-model-name) used to match a run-time override.
	ServedModelName *string `cmd:"served-model-name"`
	// Tail limits the number of most-recent log lines to show
	// (--tail), used as scrollback before following.
	Tail *int `cmd:"tail"`
	// AllSources also reads every worker/rank, not just the head
	// (--all-sources).
	AllSources bool `cmd:"all-sources"`
}

// Logs streams the log output for a running recipe (by name) or a
// specific running workload (by cluster ID).
//
// `sparkrun logs <target> [--hosts|--hosts-file|--cluster]
//
//	[--tensor-parallel] [--port] [--served-model-name] [--tail]
//	[--all-sources] --follow`
//
// The wrapper always passes --follow so the caller receives a live
// stream of new output rather than a one-shot dump (non-interactive
// callers cannot Ctrl-C a bare dump). The returned stdout [io.Reader]
// produces the log lines as sparkrun emits them; the returned stderr
// [io.Reader] produces diagnostic output. The underlying subprocess is
// terminated when the context is cancelled or when the returned kill
// function is called; callers read until io.EOF rather than closing
// the streams.
//
// io.EOF indicates a clean exit; an [ExitError] surfaces non-zero
// exits.
func (c *cliClient) Logs(ctx context.Context, target RecipeNameOrJobID, opts ...LogsOptions) (stdout, stderr io.Reader, kill func(), err error) {
	var o = list_tools.FirstOrZero(opts)

	args, err := cmdline.BuildArgs([]string{"logs", target.workloadRef()}, o, "--follow")
	if err != nil {
		return nil, nil, nil, err
	}

	return c.streamCmd(ctx, "logs", args...)
}
