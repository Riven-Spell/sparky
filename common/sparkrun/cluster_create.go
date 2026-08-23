package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
)

// ClusterCreateOptions configures [Client.ClusterCreate].
//
// sparkrun itself rejects calls without Hosts ("Error: No hosts
// provided.").
type ClusterCreateOptions struct {
	// Hosts is the replacement host list ([Hosts] or [HostsFromFile])
	// and must be non-empty -- sparkrun rejects calls without one.
	Hosts HostsList `cmd:"hosts"`
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
	// SetDefault, when true, also marks this new cluster as the default
	// (--default).
	SetDefault bool `cmd:"default"`
}

// ClusterCreate creates a new named cluster.
//
// `sparkrun cluster create <name> [--hosts] [opts...]`
func (c *cliClient) ClusterCreate(ctx context.Context, name string, opts ...ClusterCreateOptions) error {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"cluster", "create", name}, o)
	if err != nil {
		return err
	}
	return c.plainCmd(ctx, args...)
}
