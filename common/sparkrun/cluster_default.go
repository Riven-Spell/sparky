package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// ClusterDefaultOptions configures [Client.ClusterDefault]. Currently
// empty -- reserved so future options can be added without breaking
// callers.
type ClusterDefaultOptions struct{}

// ClusterDefault returns the current default cluster.
//
// `sparkrun cluster default --json`
//   - exit 0 with a JSON object: returns *ClusterSummary.
//   - exit 0 with the literal `null`: returns (nil, nil) -- this is
//     how sparkrun signals "no default is set".
func (c *cliClient) ClusterDefault(ctx context.Context, opts ...ClusterDefaultOptions) (*sparkrun_models.ClusterSummary, error) {
	_ = list_tools.FirstOrZero(opts)
	stdout, _, err := c.runCmd(ctx, "cluster", "default", "--json")
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(stdout)
	if bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	var out sparkrun_models.ClusterSummary
	if err := json.Unmarshal(trimmed, &out); err != nil {
		return nil, newParseError("cluster default", c.binaryPath, []string{"cluster", "default", "--json"}, trimmed, "", err)
	}
	return &out, nil
}
