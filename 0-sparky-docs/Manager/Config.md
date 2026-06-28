YAML file in format:

```yaml
cluster_name: default # sparkrun cluster name; will be sent to command line; `default` if not present
manager_host: "192.168.100.1:8080"
webui_host: "10.0.126.1:8080"
agent_port: 8081
models:
  - recipe: "@official/llama-3-70b"
    nodes: 2 # default 1 if not present
  - recipe: "@official/deepseek-v4-flash"
    nodes: 4
    nickname: "deepseek-v4-flash" #nickname allowed; not required. if not used, inherits recipe name.
eviction_idle_duration: "15m" # default 15m
health_check_frequency: "1m"
```

# Reload
The administrator may want to reload the configuration, to change values, to add models, etc.

If new models were added, simply add them to the list of models returned by `/v1/models`.

If models were removed, put them on the "hit list", and allow them to become evicted automatically. Do not remove them, as somebody may be using them.

# Removed models
The best decision for removing a model is to simply hard-break with advance notification. Since this is designed to be ran in a small user group (<5 users), it's not hard to notify.