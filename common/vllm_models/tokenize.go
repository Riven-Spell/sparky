package vllm_models

type TokenizeChatRequest struct {
	Model                *string                      `json:"model,omitempty"`
	Messages             []any                        `json:"messages"`
	AddGenerationPrompt  *bool                        `json:"add_generation_prompt,omitempty"`
	ReturnTokenStrs      *bool                        `json:"return_token_strs,omitempty"`
	ContinueFinalMessage *bool                        `json:"continue_final_message,omitempty"`
	AddSpecialTokens     *bool                        `json:"add_special_tokens,omitempty"`
	ChatTemplate         *string                      `json:"chat_template,omitempty"`
	ChatTemplateKwargs   map[string]any               `json:"chat_template_kwargs,omitempty"`
	Tools                []ChatCompletionToolsParam    `json:"tools,omitempty"`
}

type TokenizeCompletionRequest struct {
	Model            *string `json:"model,omitempty"`
	Prompt           string  `json:"prompt"`
	AddSpecialTokens *bool   `json:"add_special_tokens,omitempty"`
	ReturnTokenStrs  *bool   `json:"return_token_strs,omitempty"`
}

type DetokenizeRequest struct {
	Model  *string `json:"model,omitempty"`
	Tokens []int   `json:"tokens"`
}
