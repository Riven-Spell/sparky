package models

type AgentHealth struct {
	Id            string         `json:"id"`
	Alive         bool           `json:"alive"`
	MemoryUsedGb  float64        `json:"memory_used_gb"`
	Models        []ClusterModel `json:"models"`
}
