package sparkrun

import (
	"testing"
)

func TestRecipeListOptionsArgs(t *testing.T) {
	reg := "community"
	runtime := "vllm"
	got := mustArgs(t, []string{"recipe", "list"}, RecipeListOptions{
		Registry: &reg,
		Runtime:  &runtime,
		All:      true,
	})
	want := []string{
		"recipe", "list",
		"--registry", "community",
		"--runtime", "vllm",
		"--all",
	}
	assertEqualArgs(t, got, want)
}

func TestRecipeListOptionsMinimal(t *testing.T) {
	got := mustArgs(t, []string{"recipe", "list"}, RecipeListOptions{})
	assertEqualArgs(t, got, []string{"recipe", "list"})
}

func TestRecipeSearchOptionsArgs(t *testing.T) {
	runtime := "sglang"
	got := mustArgs(t, []string{"recipe", "search", "qwen"}, RecipeSearchOptions{
		Runtime: &runtime,
	})
	want := []string{"recipe", "search", "qwen", "--runtime", "sglang"}
	assertEqualArgs(t, got, want)
}

func TestRecipeShowOptionsArgs(t *testing.T) {
	tp := 2
	mem := 0.9
	got := mustArgs(t, []string{"recipe", "show", "m"}, RecipeShowOptions{
		NoVram:         true,
		TensorParallel: &tp,
		GPUMem:         &mem,
	})
	want := []string{"recipe", "show", "m", "--no-vram", "--tensor-parallel", "2", "--gpu-mem", "0.9"}
	assertEqualArgs(t, got, want)
}

func TestRecipeVramOptionsArgs(t *testing.T) {
	tp := 2
	mlen := 8192
	mem := 0.9
	got := mustArgs(t, []string{"recipe", "vram", "m"}, RecipeVramOptions{
		TensorParallel: &tp,
		MaxModelLen:    &mlen,
		GPUMem:         &mem,
		NoAutoDetect:   true,
	})
	want := []string{
		"recipe", "vram", "m",
		"--tensor-parallel", "2",
		"--max-model-len", "8192",
		"--gpu-mem", "0.9",
		"--no-auto-detect",
	}
	assertEqualArgs(t, got, want)
}
