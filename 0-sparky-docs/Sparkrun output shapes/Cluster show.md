`sparkrun cluster show <cluster_name> --json`:

```json
{
    "name": "default",
    "hosts": [
        "192.168.100.1",
        "192.168.100.2"
    ],
    "description": "",
    "user": "clusteruser",
    "default": true
}
```

Note: the `user` field is included in current sparkrun releases but
may be absent in older versions. Older docs sometimes omit it.
