package models

type AgentHealth struct {
	Alive        bool           `json:"alive"`
	MemoryUsedGb float64        `json:"memory_used_gb"`
	Models       []ClusterModel `json:"models"`
}
