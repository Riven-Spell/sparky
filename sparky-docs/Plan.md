`sparky` is comprised of an agent for each unit in the cluster, and a single "manager" that provides a webUI, and a router to live LLMs.

# Manager
## Integration
Provides an openAI compatible API for model access, and keeps track of both running and available models. If a model that isn't currently loaded is requested, it is auto loaded if there is a free node, or node(s) are candidate for eviction based upon least recent use, and a minimum time gate.

A time gate is mandatory, otherwise if more quota is required than available, models will thrash in and out. Loading a model is very time consuming, so this isn't ideal if it has to happen at all.
### Model Management
#### Loading
`sparkrun run <model> -H <targets>` actually is a lot cleaner to just run on the manager side. there's no lack of details; and it just does what we expect. We can also earmark the exact endpoint in memory at this time, but we'll have to look it up again if the manager ever crashes.
#### Eviction
In configuration, the manager should consider models "ready to evict" after a certain period of time (say, a default of 15 minutes) of disuse; i.e. no calls to inference endpoints. The manager should also have a switch for 
## WebUI
Some things inherently demand manual intervention. What if a user wants to add a new model recipe? Forcefully unload a model despite it's time-gate? [Request an alternate application is loaded?]([[Stretch Goals]]) etc.
## v1 support window

- Cluster status
	- Node availability
	- Inference-ready models
	- Log fetching
- Model management
	- Loading (allow selecting specific hosts)
	- Eviction
- Config reloading
# Agent

## Purpose
Calling `sparkrun cluster status -H <target>` only returns pertinent details like model port *provided* that it is called *on* the spark running the model. Thus, a lightweight agent must be installed on the client machine to observe this.
## Self discovery
Proactively, frequently (say, every 30s-1m), the agent will run `sparkrun cluster status -H localhost --json` to expose which models are available, and for models with a HTTP interface hosted on the local box, a no-op interface will be called to make sure the engine is still running. If it is not, the model will be evicted.
## Health checks - `GET http://<agent>/v1/health`
What models are loaded on this node? Is the node even still alive? How much memory is in use? General diagnostic information.
