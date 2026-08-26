package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// RecipeSearchOptions configures [Client.RecipeSearch].
type RecipeSearchOptions struct {
	// Registry filters results by registry name (--registry).
	Registry *string `cmd:"registry"`
	// Runtime filters results by runtime, e.g. vllm, sglang, llama-cpp
	// (--runtime).
	Runtime *string `cmd:"runtime"`
	// All includes hidden registry recipes (-a / --all).
	All bool `cmd:"all"`
}

// RecipeSearch finds recipes by name, model, or description.
//
// `sparkrun recipe search <query> [--registry] [--runtime] [-a] [--json]`
//
// Returns a (possibly empty) slice.
func (c *cliClient) RecipeSearch(ctx context.Context, query string, opts ...RecipeSearchOptions) ([]sparkrun_models.RecipeSummary, error) {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"recipe", "search", query}, o, "--json")
	if err != nil {
		return nil, err
	}
	return c.recipeSummaries(ctx, args)
}
