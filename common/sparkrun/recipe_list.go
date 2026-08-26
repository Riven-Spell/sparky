package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// RecipeListOptions configures [Client.RecipeList].
type RecipeListOptions struct {
	// Registry filters results by registry name (--registry).
	Registry *string `cmd:"registry"`
	// Runtime filters results by runtime, e.g. vllm, sglang, llama-cpp
	// (--runtime).
	Runtime *string `cmd:"runtime"`
	// All includes hidden registry recipes (-a / --all).
	All bool `cmd:"all"`
}

// RecipeList returns every available recipe from all registries.
//
// `sparkrun recipe list [query] [--registry] [--runtime] [-a] [--json]`
//
// Returns a (possibly empty) slice. An empty query lists all recipes.
func (c *cliClient) RecipeList(ctx context.Context, query string, opts ...RecipeListOptions) ([]sparkrun_models.RecipeSummary, error) {
	var o = list_tools.FirstOrZero(opts)
	prefix := []string{"recipe", "list"}
	if query != "" {
		prefix = append(prefix, query)
	}
	args, err := cmdline.BuildArgs(prefix, o, "--json")
	if err != nil {
		return nil, err
	}
	return c.recipeSummaries(ctx, args)
}
