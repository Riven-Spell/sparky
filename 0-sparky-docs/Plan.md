`sparky` is comprised of an agent for each unit in the cluster, and a single "manager" that provides a webUI, and a router to live LLMs.

# Manager
## Integration
Provides an openAI compatible API for model access. Observe [[Model Management]] for details about automagic model loading.
## WebUI
Some things inherently demand manual intervention. What if a user wants to add a new model recipe? Forcefully unload a model despite it's time-gate? [[Stretch Goals]] etc. Check [[MVP Support Window]] for details for the initial feature support targets.
# Agent

## Purpose
Calling `sparkrun cluster status -H <target>` only returns pertinent details like model port *provided* that it is called *on* the spark running the model. Thus, a lightweight agent must be installed on the client machine to observe this.
## Self discovery
Proactively, frequently (say, every 30s-1m), the agent will run `sparkrun cluster status -H localhost --json` to expose which models are available, and for models with a HTTP interface hosted on the local box, a no-op interface will be called to make sure the engine is still running. If it is not, the model will be evicted. Documented further in [[Lifecycle]]
## Health checks - `GET http://<agent>/v1/health`
What models are loaded on this node? Is the node even still alive? How much memory is in use? General diagnostic information.
