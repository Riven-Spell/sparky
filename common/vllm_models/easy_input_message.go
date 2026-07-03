package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type EasyInputMessageParam struct {
	Content any                  `json:"content"`
	Role    vllm_enum.ChatRole   `json:"role"`
	Phase   *vllm_enum.MessagePhase `json:"phase,omitempty"`
	Type    string               `json:"type"`
}
