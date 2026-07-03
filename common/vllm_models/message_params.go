package vllm_models

type ChatCompletionSystemMessageParam struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type ChatCompletionUserMessageParam struct {
	Content any    `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type ChatCompletionAssistantMessageParam struct {
	Role         string                `json:"role"`
	Content      any                   `json:"content,omitempty"`
	Name         string                `json:"name,omitempty"`
	Refusal      *string               `json:"refusal,omitempty"`
	FunctionCall *FunctionCall         `json:"function_call,omitempty"`
	ToolCalls    []ChatCompletionMessageFunctionToolCallParam `json:"tool_calls,omitempty"`
}

type ChatCompletionToolMessageParam struct {
	Content    any    `json:"content"`
	Role       string `json:"role"`
	ToolCallID string `json:"tool_call_id"`
}

type ChatCompletionFunctionMessageParam struct {
	Content *string `json:"content"`
	Name    string  `json:"name"`
	Role    string  `json:"role"`
}

type ChatCompletionDeveloperMessageParam struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type ChatCompletionMessageFunctionToolCallParam struct {
	ID       string   `json:"id"`
	Function Function `json:"function"`
	Type     string   `json:"type"`
}

type ChatCompletionMessageCustomToolCallParam struct {
	ID     string `json:"id"`
	Custom Custom `json:"custom"`
	Type   string `json:"type"`
}

type Custom struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}
