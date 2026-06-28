## Allocation
1 model takes up, bare minimum, one spark. This could change in the future, but for now it is the safest policy, unless GPU memory is managed more strictly. Startup for even small models can result in OOMs if not properly configured. Also we're splitting memory bandwidth there, so there's an extra price if multiple models run on one spark.

Some models take up more than one spark (i.e. deepseek v4 flash takes bare minimum two sparks)
## Loading
`sparkrun run <model> -H <targets>` actually is a lot cleaner to just run on the manager side. there's no lack of details; and it just does what we expect. We can also earmark the exact endpoint in memory at this time, but we'll have to look it up again if the manager ever crashes.

Run with `--label "sparky.launchtime=<UNIX timestamp>`, preserve model launch time for natively launched containers. (this shows up under the `runtime_info` field in [[Cluster status]] as `container_sparky_launchtime`)
## Eviction
In configuration, the manager should consider models "ready to evict" after a certain period of time (say, a default of 15 minutes) of disuse; i.e. no calls to inference endpoints. The manager webUI should also allow for manual eviction.
## Tracking
The manager needs to know which sparkrun IDs correspond to which models, so these should be included in the model's health check requests.

# Flows
## `/v1/chat/completions` is called against a model that is unloaded.

1) Check that the model is even registered. If not, bad request.
2) Check that the model could fit in the current cluster. i.e., are there idle nodes that could load it?
	1) If there's not enough space, can we evict anybody to make space? Follow [[#Eviction]] for details.
	2) If there's not enough space, and no room to evict, `503 Service Unavailable` [[Error Codes#CannotEvict Retriable]] is returned. `Retry-After` is the time until the next eviction would become available.
	3) If we can evict to make space, begin model eviction (via `sparkrun stop`), and progress. Let model eviction complete _first_; shouldn't take more than 5-10s. If we evict, [[Error Codes#ModelLoading Retriable]] up front, then begin long processes in background to make return quick.
3) Return [[Error Codes#ModelLoading Retriable]] with `Retry-After` set to 5 minutes
## `/v1/chat/completions` is called against a model that has been deregistered, ~~but has a substitute~~

Substituting is no longer allowed, so it instead returns [[Error Codes#ModelNotRegistered]].