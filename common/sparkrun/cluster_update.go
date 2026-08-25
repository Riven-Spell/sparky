package sparkrun

import (
	"context"

	"github.com/Riven-Spell/sparky/common/cmdline"
)

// ClusterUpdateOptions configures [Client.ClusterUpdate].
//
// The first two are mutually exclusive with the latter two; sparkrun
// rejects calls that mix them. The wrapper passes whatever the caller
// set without re-validating.
type ClusterUpdateOptions struct {
	// Hosts is a [HostsDiff] describing the host-list mutation:
	//   - [Hosts]        → replace the host list (--hosts a,b,c).
	//   - [HostsFromFile]→ replace from a file (--hosts-file path).
	//   - [AddHosts]     → append addresses (--add-host ...).
	//   - [RemoveHosts]  → remove addresses (--remove-host ...).
	Hosts HostsDiff `cmd:"hosts"`
	// Description is the cluster description (--description).
	Description *string `cmd:"description"`
	// User is the SSH username for this cluster (--user).
	User *string `cmd:"user"`
	// CacheDir is the HuggingFace cache directory (--cache-dir).
	CacheDir *string `cmd:"cache-dir"`
	// TransferMode selects the resource transfer mode (--transfer-mode):
	// auto, local, push, or delegated.
	TransferMode *TransferMode `cmd:"transfer-mode"`
	// TransferInterface selects the network interface for transfers
	// (--transfer-interface): auto (default), cx7 (InfiniBand), or
	// mgmt (management).
	TransferInterface *TransferInterface `cmd:"transfer-interface"`
	// Topology sets the CX7 topology (--topology): none (remove),
	// direct/switch (switched fabric), or ring (3-node mesh/ring).
	Topology *Topology `cmd:"topology"`
	// InferHardware, when true, SSHes into each host, detects
	// accelerators (NVIDIA/AMD/Intel/Apple) + IB, and persists per-host
	// hardware metadata (--infer-hardware).
	InferHardware bool `cmd:"infer-hardware"`
	// Executor is the default executor selector (--executor). Pass an
	// empty string to clear.
	Executor *string `cmd:"executor"`
	// ExecutorOpts are repeatable --executor-opt key=value executor
	// options. Pass empty values to emit --executor-opt key= entries.
	ExecutorOpts KeyValueOpts `cmd:"executor-opt"`
	// ClearExecutorConfig, when true, removes all executor config
	// options from the cluster (--clear-executor-config).
	ClearExecutorConfig bool `cmd:"clear-executor-config"`
	// Scheduler is the default scheduler selector (--scheduler). Pass
	// an empty string to clear.
	Scheduler *string `cmd:"scheduler"`
	// MaxGPUMemUtil caps the fraction of GPU memory usable for
	// scheduling/fit (--max-gpu-mem-util, 0.0 < x <= 1.0). Pass 0 to
	// clear.
	MaxGPUMemUtil *float64 `cmd:"max-gpu-mem-util"`
}

// ClusterUpdate mutates an existing cluster.
//
// `sparkrun cluster update <name> [--hosts|--hosts-file|--add-host|--remove-host] [opts...]`
func (c *cliClient) ClusterUpdate(ctx context.Context, name string, opts ClusterUpdateOptions) error {
	args, err := cmdline.BuildArgs([]string{"cluster", "update", name}, opts)
	if err != nil {
		return err
	}
	return c.plainCmd(ctx, args...)
}
