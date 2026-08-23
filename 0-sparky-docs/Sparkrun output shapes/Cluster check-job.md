`sparkrun cluster check-job <TARGET> --cluster <name> --json`

Exit code 0 = running (and healthy if --check-http-models), 1 = not
running or unhealthy. **JSON is emitted in both cases.** The
wrapper preserves the parsed body even on non-zero exit so callers
can read `running`/`cluster_id`/`container_statuses` from the error.

```json
{
    "running": false,
    "cluster_id": "sparkrun_3e4169f54d14",
    "healthy": null,
    "metadata": null,
    "container_statuses": {
        "sparkrun_3e4169f54d14_node_0": false,
        "sparkrun_3e4169f54d14_head": false
    },
    "hosts": ["10.0.126.1", "10.0.126.2", "10.0.126.3", "10.0.126.4"]
}
```

Field notes:
- `healthy` is `null` unless `--check-http-models` was passed; with
  it, `true`/`false` reflects the /v1/models probe.
- `metadata` is `null` unless the recipe populates it.
- `container_statuses` is keyed by container name; values are
  bools (true = up).
- TARGET is a recipe name or a `sparkrun_<hex>` cluster ID.
