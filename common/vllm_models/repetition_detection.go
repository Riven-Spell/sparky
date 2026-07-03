package vllm_models

type RepetitionDetectionParams struct {
	MaxPatternSize int `json:"max_pattern_size,omitempty"`
	MinPatternSize int `json:"min_pattern_size,omitempty"`
	MinCount       int `json:"min_count,omitempty"`
}
