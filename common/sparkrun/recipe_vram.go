package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// RecipeVramOptions configures [Client.RecipeVram].
type RecipeVramOptions struct {
	// TensorParallel overrides tensor parallelism (--tensor-parallel).
	TensorParallel *int `cmd:"tensor-parallel"`
	// MaxModelLen overrides max sequence length (--max-model-len).
	MaxModelLen *int `cmd:"max-model-len"`
	// GPUMem overrides gpu_memory_utilization, 0.0-1.0 (--gpu-mem).
	GPUMem *float64 `cmd:"gpu-mem"`
	// NoAutoDetect skips HuggingFace model auto-detection (--no-auto-detect).
	NoAutoDetect bool `cmd:"no-auto-detect"`
}

// RecipeVram estimates VRAM usage for a recipe on DGX Spark.
//
// `sparkrun recipe vram <target> [--tensor-parallel] [--max-model-len] [--gpu-mem] [--no-auto-detect] [--json]`
func (c *cliClient) RecipeVram(ctx context.Context, target RecipeNameOrFile, opts ...RecipeVramOptions) (*sparkrun_models.RecipeVramEstimate, error) {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"recipe", "vram", target.recipeRef()}, o, "--json")
	if err != nil {
		return nil, err
	}
	var out sparkrun_models.RecipeVramEstimate
	if err := c.jsonCmd(ctx, &out, args...); err != nil {
		return nil, err
	}
	return &out, nil
}
