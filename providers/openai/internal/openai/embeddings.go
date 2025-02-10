package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultEmbeddingModel = "text-embedding-ada-002"

type embeddingPayload struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingResponsePayload struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) createEmbedding(ctx context.Context, payload *embeddingPayload) (*embeddingResponsePayload, error) {
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}

	if c.Model == "" {
		payload.Model = c.EmbeddingModel
	}

	if payload.Model == "" {
		payload.Model = defaultEmbeddingModel
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	rawurl := c.buildURL("/embeddings", c.EmbeddingModel)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawurl, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("API returned unexpected status code: %d", r.StatusCode)

		var errResp errorMessage

		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("%s", msg)
		}

		return nil, fmt.Errorf("%s: %s", msg, errResp.Error.Message)
	}

	var response embeddingResponsePayload

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &response, nil
}
