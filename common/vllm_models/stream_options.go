package vllm_models

type StreamOptions struct {
	IncludeUsage         *bool `json:"include_usage,omitempty"`
	ContinuousUsageStats *bool `json:"continuous_usage_stats,omitempty"`
}
