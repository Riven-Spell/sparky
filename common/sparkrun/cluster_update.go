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
