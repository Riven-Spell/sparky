package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
)

// ClusterDeleteOptions configures [Client.ClusterDelete]. Currently
// empty -- reserved so future options (e.g. skipping --force to keep
// sparkrun's confirm prompt) can be added without breaking callers.
type ClusterDeleteOptions struct{}

// ClusterDelete removes a saved cluster.
//
// `sparkrun cluster delete <name> --force`
//
//   - exit 1: cluster does not exist.
//
// The wrapper always passes --force: a non-interactive caller has no
// way to answer sparkrun's [y/N] prompt.
func (c *cliClient) ClusterDelete(ctx context.Context, name string, opts ...ClusterDeleteOptions) error {
	_ = list_tools.FirstOrZero(opts)
	return c.plainCmd(ctx, "cluster", "delete", name, "--force")
}
