package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type ChatCompletionRequest struct {
	Messages              []any                          `json:"messages"`
	Model                 *string                        `json:"model,omitempty"`
	FrequencyPenalty      *float64                       `json:"frequency_penalty,omitempty"`
	LogitBias             map[string]float64              `json:"logit_bias,omitempty"`
	Logprobs              *bool                          `json:"logprobs,omitempty"`
	TopLogprobs           *int                           `json:"top_logprobs,omitempty"`
	MaxTokens             *int                           `json:"max_tokens,omitempty"`
	MaxCompletionTokens   *int                           `json:"max_completion_tokens,omitempty"`
	N                     *int                           `json:"n,omitempty"`
	PresencePenalty       *float64                       `json:"presence_penalty,omitempty"`
	ResponseFormat        *ResponseFormat                `json:"response_format,omitempty"`
	Seed                  *int                           `json:"seed,omitempty"`
	Stop                  any                            `json:"stop,omitempty"`
	Stream                *bool                          `json:"stream,omitempty"`
	StreamOptions         *StreamOptions                 `json:"stream_options,omitempty"`
	Temperature           *float64                       `json:"temperature,omitempty"`
	TopP                  *float64                       `json:"top_p,omitempty"`
	Tools                 []ChatCompletionToolsParam      `json:"tools,omitempty"`
	ToolChoice            any                            `json:"tool_choice,omitempty"`
	ReasoningEffort       *vllm_enum.ReasoningEffort     `json:"reasoning_effort,omitempty"`
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
	PromptLogprobs        *int                           `json:"prompt_logprobs,omitempty"`
	AllowedTokenIDs       []int                          `json:"allowed_token_ids,omitempty"`
	AddGenerationPrompt   *bool                          `json:"add_generation_prompt,omitempty"`
	ContinueFinalMessage  *bool                          `json:"continue_final_message,omitempty"`
	AddSpecialTokens      *bool                          `json:"add_special_tokens,omitempty"`
	ChatTemplate          *string                        `json:"chat_template,omitempty"`
	ChatTemplateKwargs    map[string]any                 `json:"chat_template_kwargs,omitempty"`
	RequestID             *string                        `json:"request_id,omitempty"`
	Priority              *int                           `json:"priority,omitempty"`
	StructuredOutputs     *StructuredOutputsParams       `json:"structured_outputs,omitempty"`
	RepetitionDetection   *RepetitionDetectionParams     `json:"repetition_detection,omitempty"`
	ParallelToolCalls     *bool                          `json:"parallel_tool_calls,omitempty"`
	ThinkingTokenBudget   *int                           `json:"thinking_token_budget,omitempty"`
	IncludeReasoning      *bool                          `json:"include_reasoning,omitempty"`
	Echo                  *bool                          `json:"echo,omitempty"`
	ReturnTokenIDs        *bool                          `json:"return_token_ids,omitempty"`
	BadWords              []string                       `json:"bad_words,omitempty"`
	VllmXargs             map[string]any                 `json:"vllm_xargs,omitempty"`
	CacheSalt             *string                        `json:"cache_salt,omitempty"`
	KVTransferParams      map[string]any                 `json:"kv_transfer_params,omitempty"`
}
