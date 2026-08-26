package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// ClusterListOptions configures [Client.ClusterList]. Currently empty
// -- reserved so future options can be added without breaking callers.
type ClusterListOptions struct{}

// ClusterList returns every saved cluster.
//
// `sparkrun cluster list --json`
//
// Returns a (possibly empty) slice. JSON `[]` decodes as a zero
// value of []ClusterSummary, which is the expected empty result.
func (c *cliClient) ClusterList(ctx context.Context, opts ...ClusterListOptions) ([]sparkrun_models.ClusterSummary, error) {
	_ = list_tools.FirstOrZero(opts)
	stdout, _, err := c.runCmd(ctx, "cluster", "list", "--json")
	if err != nil {
		return nil, err
	}
	var out = make([]sparkrun_models.ClusterSummary, 0)
	if len(bytes.TrimSpace(stdout)) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(stdout, &out); err != nil {
		return nil, newParseError("cluster list", c.binaryPath, []string{"cluster", "list", "--json"}, stdout, "", err)
	}
	return out, nil
}
