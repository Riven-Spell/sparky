package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type GenerateRequest struct {
	RequestID      *string           `json:"request_id,omitempty"`
	TokenIDs       []int             `json:"token_ids"`
	Features       *MultiModalFeatures `json:"features,omitempty"`
	SamplingParams SamplingParams     `json:"sampling_params"`
	Model          *string           `json:"model,omitempty"`
	Stream         *bool             `json:"stream,omitempty"`
	StreamOptions  *StreamOptions    `json:"stream_options,omitempty"`
	CacheSalt      *string           `json:"cache_salt,omitempty"`
	Priority       *int              `json:"priority,omitempty"`
	KVTransferParams map[string]any  `json:"kv_transfer_params,omitempty"`
}

type SamplingParams struct {
	N                      *int                           `json:"n,omitempty"`
	PresencePenalty        *float64                       `json:"presence_penalty,omitempty"`
	FrequencyPenalty       *float64                       `json:"frequency_penalty,omitempty"`
	RepetitionPenalty      *float64                       `json:"repetition_penalty,omitempty"`
	Temperature            *float64                       `json:"temperature,omitempty"`
	TopP                   *float64                       `json:"top_p,omitempty"`
	TopK                   *int                           `json:"top_k,omitempty"`
	MinP                   *float64                       `json:"min_p,omitempty"`
	Seed                   *int                           `json:"seed,omitempty"`
	Stop                   any                            `json:"stop,omitempty"`
	StopTokenIDs           []int                          `json:"stop_token_ids,omitempty"`
	IgnoreEOS              *bool                          `json:"ignore_eos,omitempty"`
	MaxTokens              *int                           `json:"max_tokens,omitempty"`
	MinTokens              *int                           `json:"min_tokens,omitempty"`
	Logprobs               *int                           `json:"logprobs,omitempty"`
	PromptLogprobs         *int                           `json:"prompt_logprobs,omitempty"`
	LogprobTokenIDs        []int                          `json:"logprob_token_ids,omitempty"`
	FlatLogprobs           *bool                          `json:"flat_logprobs,omitempty"`
	Detokenize             *bool                          `json:"detokenize,omitempty"`
	SkipSpecialTokens      *bool                          `json:"skip_special_tokens,omitempty"`
	SpacesBetweenSpecialTokens *bool                     `json:"spaces_between_special_tokens,omitempty"`
	IncludeStopStrInOutput *bool                         `json:"include_stop_str_in_output,omitempty"`
	OutputKind             *vllm_enum.RequestOutputKind  `json:"output_kind,omitempty"`
	SkipClone              *bool                          `json:"skip_clone,omitempty"`
	OutputTextBufferLength *int                           `json:"output_text_buffer_length,omitempty"`
	StructuredOutputs      *StructuredOutputsParams       `json:"structured_outputs,omitempty"`
	LogitBias              map[string]float64              `json:"logit_bias,omitempty"`
	AllowedTokenIDs        []int                          `json:"allowed_token_ids,omitempty"`
	BadWords               []string                       `json:"bad_words,omitempty"`
	SkipReadingPrefixCache *bool                          `json:"skip_reading_prefix_cache,omitempty"`
	ThinkingTokenBudget    *int                           `json:"thinking_token_budget,omitempty"`
	RepetitionDetection    *RepetitionDetectionParams     `json:"repetition_detection,omitempty"`
}

type MultiModalFeatures struct {
	MmHashes       map[string][]string              `json:"mm_hashes"`
	MmPlaceholders map[string][]PlaceholderRangeInfo `json:"mm_placeholders"`
	KwargsData     map[string][]*string              `json:"kwargs_data,omitempty"`
}

type PlaceholderRangeInfo struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}
