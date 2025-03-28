package types

import (
	"github.com/yottalabsai/endorphin/pkg/services/synapse"
	"time"
)

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

// Inference types
type InferenceMessageRequest struct {
	Temperature       float64       `json:"temperature" binding:"required"`
	TopP              float64       `json:"top_p" binding:"required"`
	MaxTokens         int32         `json:"max_tokens" binding:"required"`
	FrequencyPenalty  float64       `json:"frequency_penalty" binding:"required"`
	PresencePenalty   float64       `json:"presence_penalty" binding:"required"`
	RepetitionPenalty float64       `json:"repetition_penalty" binding:"required"`
	Model             string        `json:"model" binding:"required"`
	Messages          []Message     `json:"messages" binding:"required"`
	Stream            bool          `json:"stream"`
	StreamOptions     StreamOptions `json:"stream_options"`
}

type InferenceMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
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

// TextToImage types
type TextToImageRequest struct {
	Model             string  `json:"model" binding:"required"`
	Prompt            string  `json:"prompt"`
	NumInferenceSteps int32   `json:"num_inference_steps" binding:"required"`
	GuidanceScale     float64 `json:"guidance_scale"`
	LoraWeight        float64 `json:"lora_weight"`
	Seed              int32   `json:"seed"`
	Width             int32   `json:"width"`
	Height            int32   `json:"height"`
	PagScale          float64 `json:"pag_scale"`
}

type TextToImageResponse struct {
	Created int64                  `json:"created"`
	Data    []*synapse.ImageResult `json:"data"`
}

type StatusResponse struct {
	Resources Resources `json:"resources"`
	Models    Models    `json:"models"`
}

type Resources struct {
	TotalNodes int `json:"totalNodes"`
}

type Models struct {
	List []*ModelInfo `json:"list"`
}

type ModelInfo struct {
	Model string `json:"model"`
	Count int    `json:"count"`
	TPM   int    `json:"tpm"`
}
