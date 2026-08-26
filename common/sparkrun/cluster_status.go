package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// ClusterStatusOptions configures [Client.ClusterStatus].
//
// The wrapper passes whichever non-nil host selector is set to sparkrun
// and lets it resolve.
type ClusterStatusOptions struct {
	// Hosts is the list of hosts to query (--hosts).
	// At most one of Hosts / Cluster should be set.
	Hosts []string `cmd:"hosts"`
	// Cluster uses a saved cluster by name (--cluster).
	// At most one of Hosts / Cluster should be set.
	Cluster *string `cmd:"cluster"`
	// DryRun shows what would be done (--dry-run) without executing;
	// it omits --json (mutually exclusive in sparkrun), so the caller
	// receives the plain-text output and the wrapper does not parse it.
	DryRun bool `cmd:"dry-run"`
}

// ClusterStatus lists sparkrun containers running on cluster hosts.
//
// `sparkrun cluster status [--hosts|--cluster] [--dry-run] [--json]`
//
// Returns *ClusterStatusResult on --json. Per-host SSH/connect
// errors are reported in Result.Errors; they do not fail the call.
func (c *cliClient) ClusterStatus(ctx context.Context, opts ClusterStatusOptions) (*sparkrun_models.ClusterStatusResult, error) {
	args, err := cmdline.BuildArgs([]string{"cluster", "status"}, opts, "--json")
	if err != nil {
		return nil, err
	}

	if opts.DryRun {
		if err := c.plainCmd(ctx, args...); err != nil {
			return nil, err
		}
		return nil, nil
	}

	var out sparkrun_models.ClusterStatusResult
	stdout, stderrBytes, runErr := c.runCmd(ctx, args...)
	if runErr != nil {
		return nil, runErr
	}
	if len(bytes.TrimSpace(stdout)) > 0 {
		if err := json.Unmarshal(stdout, &out); err != nil {
			return nil, newParseError("cluster status", c.binaryPath, args, stdout, string(stderrBytes), err)
		}
	}
	return &out, nil
}
