`sparkrun recipe validate <target> --json`

Checks that a recipe (by name or file) is valid. `issues` is empty
when `valid` is true.

```json
{
    "recipe": "@sparkrun-transitional/qwen3-1.7b-llama-cpp",
    "valid": true,
    "issues": []
}
```

Field notes:
- An unknown recipe name yields exit 1 with a non-JSON "not found"
  message on StdError, surfaced as an [ExitError].
- The wrapper returns (nil, nil) via [ExitError] on that path -- the
  body is not present.
