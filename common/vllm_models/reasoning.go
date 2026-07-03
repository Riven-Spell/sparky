package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type Reasoning struct {
	Effort         *vllm_enum.ReasoningEffort `json:"effort,omitempty"`
	GenerateSummary *vllm_enum.Verbosity       `json:"generate_summary,omitempty"`
	Summary        *vllm_enum.Verbosity       `json:"summary,omitempty"`
}
