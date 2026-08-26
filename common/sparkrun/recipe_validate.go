package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// RecipeValidateOptions configures [Client.RecipeValidate]. Currently
// empty -- reserved so future options can be added without breaking
// callers.
type RecipeValidateOptions struct{}

// RecipeValidate checks that a recipe (by name or file) is valid.
//
// `sparkrun recipe validate <target> [--json]`
//
// A valid recipe returns (*RecipeValidateResult) with Valid set. An
// unknown recipe yields an [ExitError] whose StdError holds sparkrun's
// "not found" message.
func (c *cliClient) RecipeValidate(ctx context.Context, target RecipeNameOrFile, opts ...RecipeValidateOptions) (*sparkrun_models.RecipeValidateResult, error) {
	_ = list_tools.FirstOrZero(opts)
	var out sparkrun_models.RecipeValidateResult
	if err := c.jsonCmd(ctx, &out, "recipe", "validate", target.recipeRef()); err != nil {
		return nil, err
	}
	return &out, nil
}
