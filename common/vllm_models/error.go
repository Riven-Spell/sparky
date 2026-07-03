package vllm_models

type ErrorInfo struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    int     `json:"code"`
}

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
}
