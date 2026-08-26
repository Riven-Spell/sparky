`sparkrun recipe list [@registry] --json` and `sparkrun recipe search <query> --json`

Both return a bare JSON array. `list` with no query returns every
available recipe; `search` filters by name/model/description (and may
scope to a registry with the `@registry` shorthand). Empty result is
`[]`. The top-level aliases `sparkrun list` / `sparkrun search` share
this shape.

```json
[
  {
    "name": "@sparkrun-transitional/qwen3-1.7b-llama-cpp",
    "file": "qwen3-1.7b-llama-cpp",
    "path": "/root/.cache/sparkrun/registries/sparkrun-transitional/transitional/recipes/qwen3/qwen3-1.7b-llama-cpp.yaml",
    "model": "Qwen/Qwen3-1.7B-GGUF:Q8_0",
    "description": "",
    "runtime": "llama-cpp",
    "min_nodes": 1,
    "tp": "",
    "gpu_mem": "",
    "registry": "sparkrun-transitional"
  }
]
```

Field notes:
- `tp` and `gpu_mem` are polymorphic: sparkrun emits them as an empty
  string when unset, or as a number when a default is set (e.g. `2`,
  `0.85`). The wrapper decodes both into a [sparkrun_models.RecipeScalar].

`sparkrun recipe show <target> --json` returns a single normalized
recipe object (see "Recipe show.md"); `sparkrun recipe vram` and
`sparkrun recipe validate` return their own single-object shapes (see
"Recipe vram.md" / "Recipe validate.md").
