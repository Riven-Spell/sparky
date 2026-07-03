package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type ResponseFormat struct {
	Type       vllm_enum.ResponseFormatType `json:"type"`
	JSONSchema *JsonSchemaResponseFormat    `json:"json_schema,omitempty"`
}

type JsonSchemaResponseFormat struct {
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Schema      map[string]any         `json:"schema,omitempty"`
	Strict      *bool                  `json:"strict,omitempty"`
}
