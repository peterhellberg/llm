package openai

import "context"

// CompletionRequest is a request to complete a completion.
type CompletionRequest struct {
	Model               string   `json:"model"`
	Prompt              string   `json:"prompt"`
	Temperature         float64  `json:"temperature"`
	MaxCompletionTokens int      `json:"max_completion_tokens,omitempty"`
	N                   int      `json:"n,omitempty"`
	FrequencyPenalty    float64  `json:"frequency_penalty,omitempty"`
	PresencePenalty     float64  `json:"presence_penalty,omitempty"`
	TopP                float64  `json:"top_p,omitempty"`
	StopWords           []string `json:"stop,omitempty"`
	Seed                int      `json:"seed,omitempty"`

	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
}

type CompletionResponse struct {
	ID      string  `json:"id,omitempty"`
	Created float64 `json:"created,omitempty"`
	Choices []struct {
		FinishReason string      `json:"finish_reason,omitempty"`
		Index        float64     `json:"index,omitempty"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		Text         string      `json:"text,omitempty"`
	} `json:"choices,omitempty"`
	Model  string `json:"model,omitempty"`
	Object string `json:"object,omitempty"`
	Usage  struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

type errorMessage struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (c *Client) setCompletionDefaults(payload *CompletionRequest) {
	if len(payload.StopWords) == 0 {
		payload.StopWords = nil
	}

	switch {
	case payload.Model != "":
	case c.Model != "":
		payload.Model = c.Model
	default:
		payload.Model = defaultChatModel
	}
}

func (c *Client) createCompletion(ctx context.Context, payload *CompletionRequest) (*ChatCompletionResponse, error) {
	c.setCompletionDefaults(payload)

	req := &ChatRequest{
		Model: payload.Model,
		Messages: []*ChatMessage{
			{Role: "user", Content: payload.Prompt},
		},
		Temperature:         payload.Temperature,
		TopP:                payload.TopP,
		MaxCompletionTokens: payload.MaxCompletionTokens,
		N:                   payload.N,
		StopWords:           payload.StopWords,
		FrequencyPenalty:    payload.FrequencyPenalty,
		PresencePenalty:     payload.PresencePenalty,
		StreamingFunc:       payload.StreamingFunc,
		Seed:                payload.Seed,
	}

	return c.createChat(ctx, req)
}
