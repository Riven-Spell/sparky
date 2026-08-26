`sparkrun recipe show <target> --json`

Returns the fully-normalized recipe for a recipe name or a recipe-YAML
file path. `metadata` and `defaults` are free-form per-runtime maps,
typed as `any` in the wrapper so callers decode only what they use.

```json
{
    "recipe_version": "2",
    "model": "Qwen/Qwen3-1.7B-GGUF:Q8_0",
    "runtime": "llama-cpp",
    "max_nodes": 1,
    "container": "ghcr.io/spark-arena/dgx-llama-cpp:latest",
    "metadata": {
        "description": "Qwen3 1.7B (Q8_0 GGUF)",
        "maintainer": "scitrera.ai <open-source-team@scitrera.com>",
        "model_params": "1.7B",
        "model_dtype": "q8_0"
    },
    "defaults": {
        "port": 8000,
        "host": "0.0.0.0",
        "n_gpu_layers": 99,
        "ctx_size": 8192
    },
    "command": "llama-server   -hf {model}   --host {host} ..."
}
```

Field notes:
- `recipe_version`, `max_nodes`, `--no-vram`/`--tensor-parallel`/`--gpu-mem`
  overrides apply here.
- The wrapper type is [sparkrun_models.RecipeDetail]. An unknown recipe
  name yields an exit-1 [ExitError] with "not found" on StdError.
