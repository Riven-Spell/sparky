package sparkrun

import (
	"context"
	"io"
	"strconv"
	"strings"
)

// ClusterMonitorOptions configures [Client.ClusterMonitor].
//
// The wrapper returns the raw io.Reader streams regardless of which
// mode was selected; the caller decides how to render them.
type ClusterMonitorOptions struct {
	// Hosts is the list of host addresses to monitor (--hosts).
	// Mutually exclusive with Cluster.
	Hosts []string
	// Cluster uses a saved cluster by name to select its hosts
	// (--cluster). Mutually exclusive with Hosts.
	Cluster *string
	// DryRun shows what would be done (--dry-run) without executing.
	DryRun bool
	// Interval is the sampling interval in seconds (--interval).
	Interval *int
	// Simple uses plain-text output (--simple) instead of the
	// interactive TUI.
	Simple bool
}

// ClusterMonitor streams per-host CPU/RAM/GPU metrics.
//
// `sparkrun cluster monitor [--hosts|--cluster]
//
//	[--dry-run] [--interval] [--simple]`
//
// The returned stdout [io.Reader] produces NDJSON: one JSON object per line.
// The returned stderr [io.Reader] produces diagnostic output from sparkrun.
// The underlying subprocess is terminated when the context is cancelled or
// when the returned kill function is called; callers read until io.EOF
// rather than closing the streams.
//
// io.EOF indicates a clean exit; an [ExitError] surfaces non-zero exits.
//
// The wrapper does not validate the mutual exclusivity of --json,
// --simple, and --dry-run at this layer -- it simply forwards
// whichever flags the caller set. sparkrun itself rejects
// combinations.
func (c *cliClient) ClusterMonitor(ctx context.Context, opts ClusterMonitorOptions) (stdout, stderr io.Reader, kill func(), err error) {
	args := []string{"cluster", "monitor"}
	if len(opts.Hosts) > 0 {
		args = append(args, "--hosts", strings.Join(opts.Hosts, ","))
	}
	if opts.Cluster != nil {
		args = append(args, "--cluster", *opts.Cluster)
	}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	if opts.Interval != nil {
		args = append(args, "--interval", strconv.Itoa(*opts.Interval))
	}
	if opts.Simple {
		args = append(args, "--simple")
	}
	// --json is the default streaming mode; --dry-run and --simple
	// override it. We deliberately don't pass --json when the
	// caller asked for --simple or --dry-run.
	if !opts.DryRun && !opts.Simple {
		args = append(args, "--json")
	}

	return c.streamCmd(ctx, "cluster monitor", args...)
}
