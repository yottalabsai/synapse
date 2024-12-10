package types

import "time"

type CreateServerlessResourceRequest struct {
	EndpointId string `json:"endpointId"    binding:"required"`
	Model      string `json:"model"   binding:"required"`
}

type ServerlessResourceResponse struct {
	ID         int64     `json:"id"`
	EndpointId string    `json:"endpointId"`
	Model      string    `json:"model"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type InferenceMessageRequest struct {
	Messages []InferenceMessage `json:"messages"`
}

type InferenceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	Logprobs     *Logprobs `json:"logprobs"`
	FinishReason string    `json:"finish_reason"`
	StopReason   *string   `json:"stop_reason"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	// Define fields for ToolCall if needed
}

type Logprobs struct {
	// Define fields for Logprobs if needed
}

type Usage struct {
	PromptTokens     int       `json:"prompt_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	PromptLogprobs   *Logprobs `json:"prompt_logprobs"`
}
