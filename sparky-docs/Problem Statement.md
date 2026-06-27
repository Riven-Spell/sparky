The DGX Spark, compared to other options for graphics compute, provides a sizeable amount of VRAM for it's price point (and especially on power usage!) at 4,000$.

This comes with catches though... this "VRAM" is actually contiguous system memory, shared between the GPU and CPU. It's relatively slow compared to a "real" gaming or server GPU. (Only about 300GB/s!)

Moreover, because they're individual nodes, and not individual GPUs, some magic is needed to manage language models across the entire cluster. No current solution provides the level of flexibility wanted in a multi-user cluster.

## Explored Solutions
### Llama-swap

| Does...                                                                   | But...                                                           |
| ------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| Makes recipes reproducible                                                | Shared runtime means all changes affect everybody                |
| Require literally no new bespoke solution; can use all existing software. | No management of cross-node resources.                           |
| Allows for many runtimes                                                  | Kind of evidently jank and intended for a single node, not many. |
|                                                                           | llama.cpp already solves this problem, amongst other things.     |
### Automatically invoking vLLM or sgLang with a single Ray distributed backend 

| Does...                                                       | But...                                                         |
| ------------------------------------------------------------- | -------------------------------------------------------------- |
| Solves the issue with GPU distribution, handily!              | Is less directly controllable which nodes the backend goes to. |
| There's less overhead, since we're not worrying about Docker! | Runtime is global. Any update affects all nodes and models.    |
| Very low complexity, like llama-swap                          | Wastes an entire useful repository of sparkrun recipes         |
|                                                               | Decreases repeatability across the board                       |
### Invoking sparkrun recipes across nodes, and worrying about distribution ourselves. (The intended solution)

| Does...                                                          | But...                                                                                                                   |
| ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| Fix the reproduceability issue                                   | Adds multiple software dependencies                                                                                      |
| No manual reproduction of Sparkrun recipes                       | Requires significantly more "manual" management of resources, increasing complexity                                      |
| Clearly define the cluster, easily!                              | Status obtaining isn't trivial; it must be done on the node that is running a model to discover where the endpoint lives |
| Provide a reliable, and easily manipulatable runtime environment | Relies upon Docker to do so, which adds a layer of complexity.                                                           |

