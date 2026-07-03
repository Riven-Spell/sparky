package vllm_models

type HTTPValidationError struct {
	Detail []ValidationError `json:"detail,omitempty"`
}

type ValidationError struct {
	Loc  []any  `json:"loc"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}
