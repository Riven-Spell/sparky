package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
)

// ClusterUnsetDefaultOptions configures [Client.ClusterUnsetDefault].
// Currently empty -- reserved so future options can be added without
// breaking callers.
type ClusterUnsetDefaultOptions struct{}

// ClusterUnsetDefault clears the default cluster setting. It is
// idempotent -- sparkrun exits 0 whether or not a default was set.
//
// `sparkrun cluster unset-default`
func (c *cliClient) ClusterUnsetDefault(ctx context.Context, opts ...ClusterUnsetDefaultOptions) error {
	_ = list_tools.FirstOrZero(opts)
	return c.plainCmd(ctx, "cluster", "unset-default")
}
