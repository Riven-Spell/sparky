package vllm_models

type StructuredOutputsParams struct {
	JSON                        *string `json:"json,omitempty"`
	Regex                       *string `json:"regex,omitempty"`
	Choice                      []string `json:"choice,omitempty"`
	Grammar                     *string `json:"grammar,omitempty"`
	JSONObject                  *bool   `json:"json_object,omitempty"`
	DisableAnyWhitespace        bool    `json:"disable_any_whitespace,omitempty"`
	DisableAdditionalProperties bool    `json:"disable_additional_properties,omitempty"`
	WhitespacePattern           *string `json:"whitespace_pattern,omitempty"`
	StructuralTag               *string `json:"structural_tag,omitempty"`
	Backend                     *string `json:"_backend,omitempty"`
	BackendWasAuto              bool    `json:"_backend_was_auto,omitempty"`
}
