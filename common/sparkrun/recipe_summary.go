package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// recipeSummaries runs a command that yields a bare array of
// [sparkrun_models.RecipeSummary] and decodes it, guaranteeing a
// non-nil empty slice when sparkrun emits `[]`. Shared by
// [Client.RecipeList] and [Client.RecipeSearch].
func (c *cliClient) recipeSummaries(ctx context.Context, args []string) ([]sparkrun_models.RecipeSummary, error) {
	stdout, _, err := c.runCmd(ctx, args...)
	if err != nil {
		return nil, err
	}
	var out = make([]sparkrun_models.RecipeSummary, 0)
	if len(bytes.TrimSpace(stdout)) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(stdout, &out); err != nil {
		return nil, newParseError(subcommandPath(args), c.binaryPath, args, stdout, "", err)
	}
	return out, nil
}
