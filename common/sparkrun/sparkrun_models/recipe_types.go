package sparkrun_models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// RecipeScalar is a value sparkrun emits as either a JSON string
// (e.g. "" for unset) or a JSON number (e.g. 2 or 0.85) in the tp and
// gpu_mem fields of recipe list entries. It decodes to a float64;
// the empty-string unset form is represented as NaN so callers can
// distinguish "not set" from zero with math.IsNaN. Use IsSet to test.
type RecipeScalar float64

// UnmarshalJSON decodes a JSON string or number into a float64.
// The empty string (unset) decodes to NaN.
func (s *RecipeScalar) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		if str == "" {
			*s = RecipeScalar(math.NaN())
			return nil
		}
		val, convErr := strconv.ParseFloat(str, 64)
		if convErr != nil {
			return fmt.Errorf("recipe scalar: %w", convErr)
		}
		*s = RecipeScalar(val)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(b, &num); err == nil {
		val, convErr := num.Float64()
		if convErr != nil {
			return fmt.Errorf("recipe scalar: %w", convErr)
		}
		*s = RecipeScalar(val)
		return nil
	}
	return fmt.Errorf("recipe scalar: cannot decode %s", string(b))
}

// IsSet reports whether the scalar was present (not the empty string).
func (s RecipeScalar) IsSet() bool {
	return !math.IsNaN(float64(s))
}

// RecipeSummary is one entry in the output of `sparkrun recipe list`
// and `sparkrun recipe search` (both emit the same bare array shape),
// and of the top-level `list`/`search` aliases.
type RecipeSummary struct {
	Name        string       `json:"name"`
	File        string       `json:"file"`
	Path        string       `json:"path"`
	Model       string       `json:"model"`
	Description string       `json:"description"`
	Runtime     string       `json:"runtime"`
	MinNodes    int          `json:"min_nodes"`
	TP          RecipeScalar `json:"tp"`
	GPUMem      RecipeScalar `json:"gpu_mem"`
	Registry    string       `json:"registry"`
}

// RecipeDetail is the body of `sparkrun recipe show --json`: the
// fully-normalized recipe. Metadata and Defaults are free-form
// per-runtime key/value maps, typed as any so callers decode the
// fields they care about without baking in every runtime's shape.
type RecipeDetail struct {
	RecipeVersion string         `json:"recipe_version"`
	Model         string         `json:"model"`
	Runtime       string         `json:"runtime"`
	MaxNodes      int            `json:"max_nodes"`
	Container     string         `json:"container"`
	Metadata      map[string]any `json:"metadata"`
	Defaults      map[string]any `json:"defaults"`
	Command       string         `json:"command"`
}

// RecipeVramEstimate is the body of `sparkrun recipe vram --json`.
// Many fields are nullable because they depend on model architecture
// metadata that may not be available (e.g. KV cache sizing for a
// recipe whose arch info is missing).
type RecipeVramEstimate struct {
	Recipe               string   `json:"recipe"`
	Model                string   `json:"model"`
	Runtime              string   `json:"runtime"`
	ModelWeightsGB       float64  `json:"model_weights_gb"`
	KVCachePerTokenBytes *float64 `json:"kv_cache_per_token_bytes"`
	KVCacheTotalGB       *float64 `json:"kv_cache_total_gb"`
	TotalPerGPUGB        float64  `json:"total_per_gpu_gb"`
	MaxModelLen          *int     `json:"max_model_len"`
	TensorParallel       int      `json:"tensor_parallel"`
	PipelineParallel     int      `json:"pipeline_parallel"`
	Warnings             []string `json:"warnings"`
	ModelParams          int64    `json:"model_params"`
	ModelDtype           string   `json:"model_dtype"`
	KVDtype              *string  `json:"kv_dtype"`
	NumLayers            *int     `json:"num_layers"`
	NumKVHeads           *int     `json:"num_kv_heads"`
	HeadDim              *int     `json:"head_dim"`
	KVArch               string   `json:"kv_arch"`
	KVArchLabel          *string  `json:"kv_arch_label"`
	KVCacheReplicated    bool     `json:"kv_cache_replicated"`
	KVEstimateIsFloor    bool     `json:"kv_estimate_is_floor"`
	GPUMemoryUtilization *float64 `json:"gpu_memory_utilization"`
	TotalGPUMemoryGB     float64  `json:"total_gpu_memory_gb"`
	UsableGPUMemoryGB    *float64 `json:"usable_gpu_memory_gb"`
	AvailableKVGB        *float64 `json:"available_kv_gb"`
	MaxContextTokens     *int     `json:"max_context_tokens"`
	ContextMultiplier    *float64 `json:"context_multiplier"`
	FitsDGXSpark         bool     `json:"fits_dgx_spark"`
}

// RecipeValidateResult is the body of `sparkrun recipe validate --json`.
// Issues is empty when Valid is true.
type RecipeValidateResult struct {
	Recipe string   `json:"recipe"`
	Valid  bool     `json:"valid"`
	Issues []string `json:"issues"`
}
