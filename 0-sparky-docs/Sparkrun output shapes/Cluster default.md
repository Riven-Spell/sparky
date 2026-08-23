`sparkrun cluster default --json`:

```json
{
    "name": "default",
    "hosts": ["10.0.126.1", "10.0.126.2", "10.0.126.3", "10.0.126.4"],
    "description": "",
    "user": "clusteruser",
    "default": true
}
```

If no default is set, sparkrun returns the literal JSON `null` on
stdout and exit code 0. The wrapper surfaces this as `(nil, nil)`.

The shape itself matches `cluster show` -- see Cluster show.md.
