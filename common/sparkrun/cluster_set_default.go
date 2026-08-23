package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
)

// ClusterSetDefaultOptions configures [Client.ClusterSetDefault].
// Currently empty -- reserved so future options can be added without
// breaking callers.
type ClusterSetDefaultOptions struct{}

// ClusterSetDefault marks a saved cluster as the default. Exit 1
// indicates the named cluster does not exist.
//
// `sparkrun cluster set-default <name>`
func (c *cliClient) ClusterSetDefault(ctx context.Context, name string, opts ...ClusterSetDefaultOptions) error {
	_ = list_tools.FirstOrZero(opts)
	return c.plainCmd(ctx, "cluster", "set-default", name)
}
