package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
)

// ClusterShowOptions configures [Client.ClusterShow]. Currently empty
// -- reserved so future options can be added without breaking callers.
type ClusterShowOptions struct{}

// ClusterShow returns the saved definition of a named cluster.
//
// `sparkrun cluster show <name> --json`
//
//   - exit 0: returns *ClusterSummary.
//   - exit 1: the named cluster does not exist; the wrapper returns
//     an [ExitError] whose StdError holds the human-readable message
//     ("Error: Cluster ... not found").
func (c *cliClient) ClusterShow(ctx context.Context, name string, opts ...ClusterShowOptions) (*ClusterSummary, error) {
	_ = list_tools.FirstOrZero(opts)
	var out ClusterSummary
	if err := c.jsonCmd(ctx, &out, "cluster", "show", name); err != nil {
		return nil, err
	}
	return &out, nil
}
