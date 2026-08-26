`sparkrun recipe vram <target> --json`

Estimates VRAM usage for a recipe on a DGX Spark. Many fields are
null/absent when model architecture metadata is missing (KV cache
sizing, max context, etc.), so the wrapper types them as pointers.


```json
{
    "recipe": "@sparkrun-transitional/qwen3-1.7b-llama-cpp",
    "model": "Qwen/Qwen3-1.7B-GGUF:Q8_0",
    "runtime": "llama-cpp",
    "model_weights_gb": 1.682201400399208,
    "kv_cache_per_token_bytes": null,
    "kv_cache_total_gb": null,
    "total_per_gpu_gb": 1.682201400399208,
    "max_model_len": null,
    "tensor_parallel": 1,
    "pipeline_parallel": 1,
    "warnings": ["Missing architecture info"],
    "model_params": 1700000000,
    "model_dtype": "q8_0",
    "kv_dtype": null,
    "num_layers": null,
    "num_kv_heads": null,
    "head_dim": null,
    "kv_arch": "dense",
    "kv_arch_label": null,
    "kv_cache_replicated": false,
    "kv_estimate_is_floor": false,
    "gpu_memory_utilization": null,
    "total_gpu_memory_gb": 121.0,
    "usable_gpu_memory_gb": null,
    "available_kv_gb": null,
    "max_context_tokens": null,
    "context_multiplier": null,
    "fits_dgx_spark": true
}
```

Field notes:
- `fits_dgx_spark` is the overall verdict; `warnings` explain gaps.
- The wrapper type is [sparkrun_models.RecipeVramEstimate].
