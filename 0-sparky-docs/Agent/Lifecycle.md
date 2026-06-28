# Boot
1) Run [[Cluster status]] targeting localhost to discover what models we are hosting. 
2) Fire up the HTTP(S) server for the agent.
3) Ping the manager with it's health data, retrying a few times if it isn't quite up yet, to make ourselves immediately available.
# Self-Monitoring
Every 5s, trigger [[Cluster status]] targeting localhost, which should execute fairly immediately, since it only needs to poke the local docker daemon to get this info. 
## Model crash observation
Additionally, call `/v1/models` on the hosts observed to check if the model's server is still live. If it isn't, it's a good indicator that the model is no longer alive. 

Should the model's server not be listening for several health checks, evict the model and inform the manager that it has crashed.

## Model sideload observation
If a model we didn't see before is now present, it maybe got sideloaded! This gets reported to the manager with a sudden agent health push.

The new model becomes available for inference. Check [[Cluster status]] runtime info for detail:
- `container_sparky_nolist` (`--label sparky.nolist=true`) - if set, this model takes up the node, but isn't listed. Treat the node as unavailable for another model.
- `container_sparky_noevict` (`--label sparky.noevict=true`) - if set, this model does not auto-evict. manual eviction is required.

Note to self: try to write these values with `--label` on `sparkrun run`. this is output in [[Cluster status]]. this allows for us to set these with a real model without writing a custom yaml file.

If the model is to be listed, treat it like any other model. List it in `/v1/models`.
If the model is to be evicted, treat it like any other model. Evict it after a certain amount of time.
# Health Checks
Return the details from the last most recent health check, unless the `fresh=true` query parameter is set.