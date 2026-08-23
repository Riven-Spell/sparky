package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
)

// ClusterCheckJobOptions configures [Client.ClusterCheckJob].
//
// Hosts / Cluster select which hosts are queried.
// TensorParallel, Port, and ServedModelName override how sparkrun
// derives a cluster ID from a recipe name (so re-checks against
// the same recipe with the same overrides hit the same workload).
// CheckHTTPModels additionally probes /v1/models on the inference
// server to populate Healthy.
type ClusterCheckJobOptions struct {
	// Hosts is the host list to check against (--hosts or --hosts-file).
	Hosts HostsList `cmd:"hosts"`
	// Cluster uses a saved cluster by name (--cluster) to resolve the
	// host set.
	Cluster *string `cmd:"cluster"`
	// TensorParallel is the tensor-parallelism override (--tensor-parallel),
	// used for cluster_id generation so re-checks
	// against the same recipe with the same overrides hit the same
	// workload.
	TensorParallel *int `cmd:"tensor-parallel"`
	// Port is the inference-server port override (--port), used for
	// cluster_id generation and the health check.
	Port *int `cmd:"port"`
	// ServedModelName is the served-model-name override
	// (--served-model-name), used for cluster_id generation.
	ServedModelName *string `cmd:"served-model-name"`
	// CheckHTTPModels also verifies that the inference server responds
	// to health checks at /v1/models (--check-http-models) to populate
	// Healthy.
	CheckHTTPModels *bool `cmd:"check-http-models"`
}

// ClusterCheckJob reports whether a recipe (by name) or a running
// workload (by cluster ID) is up. target is a [RecipeNameOrJobID]
// because sparkrun accepts either form.
//
// `sparkrun cluster check-job <target> [--hosts|--cluster]
//
//	[--tensor-parallel] [--port] [--served-model-name] [--check-http-models] [--json]`
//
// Return semantics (per design):
//   - exit 0: returns (*ClusterCheckJobResult, nil).
//   - exit != 0 with parseable JSON: returns (nil, *ExitError) with
//     the parsed body on Result().
//   - exit != 0 with non-JSON body: returns (nil, *ExitError) with
//     the raw bytes on Raw().
//
// Callers that don't care about the body can simply treat any
// non-nil error as "not running"; callers that want the body
// should type-assert [ExitError] and read Result() / Raw().
func (c *cliClient) ClusterCheckJob(ctx context.Context, target RecipeNameOrJobID, opts ...ClusterCheckJobOptions) (*ClusterCheckJobResult, error) {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"cluster", "check-job", target.workloadRef()}, o, "--json")
	if err != nil {
		return nil, err
	}

	stdout, stderrBytes, runErr := c.runCmd(ctx, args...)
	trimmed := bytes.TrimSpace(stdout)

	if runErr == nil {
		if len(trimmed) == 0 {
			return nil, nil
		}
		var out ClusterCheckJobResult
		if err := json.Unmarshal(trimmed, &out); err != nil {
			return nil, newParseError("cluster check-job", c.binaryPath, args, trimmed, string(stderrBytes), err)
		}
		return &out, nil
	}

	// Non-zero exit. The wrapper preserves any decoded body so
	// callers that want "what would running look like?" can read it
	// from the error. Only an *ExitError carries the body; exec and
	// usage failures are passed through as-is.
	var xe *ExitError
	if !errors.As(runErr, &xe) {
		return nil, runErr
	}
	if len(trimmed) > 0 {
		var parsed ClusterCheckJobResult
		if err := json.Unmarshal(trimmed, &parsed); err == nil {
			xe.result = &parsed
		} else {
			xe.raw = append([]byte(nil), trimmed...)
		}
	}
	return nil, xe
}
