package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type CompletionRequest struct {
	Model                 *string                        `json:"model,omitempty"`
	Prompt                any                            `json:"prompt,omitempty"`
	Echo                  *bool                          `json:"echo,omitempty"`
	FrequencyPenalty      *float64                       `json:"frequency_penalty,omitempty"`
	LogitBias             map[string]float64              `json:"logit_bias,omitempty"`
	Logprobs              *int                           `json:"logprobs,omitempty"`
	MaxTokens             *int                           `json:"max_tokens,omitempty"`
	N                     *int                           `json:"n,omitempty"`
	PresencePenalty       *float64                       `json:"presence_penalty,omitempty"`
	Seed                  *int                           `json:"seed,omitempty"`
	Stop                  any                            `json:"stop,omitempty"`
	Stream                *bool                          `json:"stream,omitempty"`
	StreamOptions         *StreamOptions                 `json:"stream_options,omitempty"`
	Suffix                *string                        `json:"suffix,omitempty"`
	Temperature           *float64                       `json:"temperature,omitempty"`
	TopP                  *float64                       `json:"top_p,omitempty"`
	User                  *string                        `json:"user,omitempty"`
	UseBeamSearch         *bool                          `json:"use_beam_search,omitempty"`
	TopK                  *int                           `json:"top_k,omitempty"`
	MinP                  *float64                       `json:"min_p,omitempty"`
	RepetitionPenalty     *float64                       `json:"repetition_penalty,omitempty"`
	LengthPenalty         *float64                       `json:"length_penalty,omitempty"`
	StopTokenIDs          []int                          `json:"stop_token_ids,omitempty"`
	IncludeStopStrInOutput *bool                         `json:"include_stop_str_in_output,omitempty"`
	IgnoreEOS             *bool                          `json:"ignore_eos,omitempty"`
	MinTokens             *int                           `json:"min_tokens,omitempty"`
	SkipSpecialTokens     *bool                          `json:"skip_special_tokens,omitempty"`
	SpacesBetweenSpecialTokens *bool                     `json:"spaces_between_special_tokens,omitempty"`
	TruncatePromptTokens  *int                           `json:"truncate_prompt_tokens,omitempty"`
	TruncationSide        *vllm_enum.TruncationSide      `json:"truncation_side,omitempty"`
	AllowedTokenIDs       []int                          `json:"allowed_token_ids,omitempty"`
	PromptLogprobs        *int                           `json:"prompt_logprobs,omitempty"`
	PromptEmbeds          any                            `json:"prompt_embeds,omitempty"`
	AddSpecialTokens      *bool                          `json:"add_special_tokens,omitempty"`
	ResponseFormat        *ResponseFormat                `json:"response_format,omitempty"`
	StructuredOutputs     *StructuredOutputsParams       `json:"structured_outputs,omitempty"`
	Priority              *int                           `json:"priority,omitempty"`
	RequestID             *string                        `json:"request_id,omitempty"`
	ReturnTokensAsTokenIDs *bool                         `json:"return_tokens_as_token_ids,omitempty"`
	ReturnTokenIDs        *bool                          `json:"return_token_ids,omitempty"`
	CacheSalt             *string                        `json:"cache_salt,omitempty"`
	KVTransferParams      map[string]any                 `json:"kv_transfer_params,omitempty"`
	VllmXargs             map[string]any                 `json:"vllm_xargs,omitempty"`
	RepetitionDetection   *RepetitionDetectionParams     `json:"repetition_detection,omitempty"`
	ThinkingTokenBudget   *int                           `json:"thinking_token_budget,omitempty"`
}
