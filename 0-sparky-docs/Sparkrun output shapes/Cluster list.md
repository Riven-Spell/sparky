`sparkrun cluster list --json`:

```json
[
    {
        "name": "default",
        "hosts": ["10.0.126.1", "10.0.126.2", "10.0.126.3", "10.0.126.4"],
        "description": "",
        "user": "clusteruser",
        "default": true
    },
    {
        "name": "pair-a",
        "hosts": ["10.0.126.1", "10.0.126.2"],
        "description": "",
        "user": "clusteruser",
        "default": false
    }
]
```

Note: a **bare JSON array** is returned, not a wrapper object. The
element shape matches `cluster show`. An empty list is `[]`.

This is also the same shape (as a single object) used by
`cluster default` -- see Cluster default.md.
