package vllm_models

type Function struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

type FunctionCall struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

type FunctionTool struct {
	DeferLoading *bool           `json:"defer_loading,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Name         string          `json:"name"`
	Parameters   map[string]any  `json:"parameters,omitempty"`
	Strict       *bool           `json:"strict,omitempty"`
	Type         string          `json:"type"`
}

type ChatCompletionNamedFunction struct {
	Name string `json:"name"`
}

type ChatCompletionNamedToolChoiceParam struct {
	Function ChatCompletionNamedFunction `json:"function"`
	Type     string                      `json:"type"`
}

type ChatCompletionToolsParam struct {
	DeferLoading *bool          `json:"defer_loading,omitempty"`
	Function     FunctionTool   `json:"function"`
	Type         string         `json:"type"`
}
