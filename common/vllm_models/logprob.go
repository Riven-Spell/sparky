package vllm_models

type Logprob struct {
	Token       string           `json:"token"`
	Bytes       []int            `json:"bytes"`
	Logprob     float64          `json:"logprob"`
	TopLogprobs []LogprobTopLogprob `json:"top_logprobs"`
}

type LogprobTopLogprob struct {
	Token   string  `json:"token"`
	Bytes   []int   `json:"bytes"`
	Logprob float64 `json:"logprob"`
}
