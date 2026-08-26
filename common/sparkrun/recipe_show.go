package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// RecipeShowOptions configures [Client.RecipeShow].
type RecipeShowOptions struct {
	// NoVram skips the VRAM estimation that sparkrun otherwise attaches
	// to the normalized recipe (--no-vram).
	NoVram bool `cmd:"no-vram"`
	// TensorParallel overrides tensor parallelism (--tensor-parallel).
	TensorParallel *int `cmd:"tensor-parallel"`
	// GPUMem overrides GPU memory utilization, 0.0-1.0 (--gpu-mem).
	GPUMem *float64 `cmd:"gpu-mem"`
}

// RecipeShow returns the fully-normalized recipe for a recipe name or
// a recipe-YAML file.
//
// `sparkrun recipe show <target> [--no-vram] [--tensor-parallel] [--gpu-mem] [--json]`
func (c *cliClient) RecipeShow(ctx context.Context, target RecipeNameOrFile, opts ...RecipeShowOptions) (*sparkrun_models.RecipeDetail, error) {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"recipe", "show", target.recipeRef()}, o, "--json")
	if err != nil {
		return nil, err
	}
	var out sparkrun_models.RecipeDetail
	if err := c.jsonCmd(ctx, &out, args...); err != nil {
		return nil, err
	}
	return &out, nil
}
