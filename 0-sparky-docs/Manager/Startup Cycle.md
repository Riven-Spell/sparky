On bootup, scan through [[Cluster show]] and try say hi to agents on `agent_port` for each host. This establishes live and dead agents. If an agent wasn't up yet, it will probably be up by the next health check window. This produces a two way duplicate call, _yes_, but, this handles crashes on both ends.

Agents may not all be online on manager startup, so we poll at an accelerated rate (5s) until at least one agent is online. This comes with a maximum timeout of 3 repeats of `health_check_frequency`.

Until enough agents are online, the server will respond to inference requests with `503 - Service Unavailable` ([ModelTooLarge]([[Error Codes#ModelTooLarge]])) until enough resources are available for the requested model.