package types

type Common struct {
	ClientId    string `json:"name"`
	Timestamp   int64  `json:"timestamp"`
	MessageId   string `json:"message_id"`
	MessageType string `json:"message_type"`
}

const (
	Error       string = "error"
	Ping        string = "ping"
	Inference   string = "inference"
	TextToImage string = "text_to_image"
)

type PingRequest struct {
	Common
}

type PingResponse struct {
	Common
}

type InferenceRequest struct {
	Temperature       float64        `json:"temperature"`
	TopP              float64        `json:"top_p"`
	MaxTokens         int32          `json:"max_tokens"`
	FrequencyPenalty  float64        `json:"frequency_penalty"`
	PresencePenalty   float64        `json:"presence_penalty"`
	RepetitionPenalty float64        `json:"repetition_penalty"`
	Model             string         `json:"model"`
	Messages          []*Message     `json:"messages"`
	Stream            bool           `json:"stream"`
	StreamOptions     *StreamOptions `json:"stream_options"`
}

type InferenceResult struct {
	Common
	Message string `json:"message"`
}

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
