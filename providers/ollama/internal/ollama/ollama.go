package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

const maxBufferSize = 512 * 1000

var ErrNoURL = fmt.Errorf("no url provided")

type Client struct {
	base *url.URL
	http *http.Client
}

func NewClient(base *url.URL, client *http.Client) (*Client, error) {
	if base == nil {
		return nil, ErrNoURL
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}

	return &Client{
		base: base,
		http: client,
	}, nil
}

type GenerateResponseFunc func(GenerateResponse) error

func (c *Client) Generate(ctx context.Context, req *GenerateRequest, fn GenerateResponseFunc) error {
	return c.stream(ctx, http.MethodPost, "/api/generate", req, func(data []byte) error {
		var resp GenerateResponse

		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		return fn(resp)
	})
}

type ChatResponseFunc func(ChatResponse) error

func (c *Client) GenerateChat(ctx context.Context, req *ChatRequest, fn ChatResponseFunc) error {
	return c.stream(ctx, http.MethodPost, "/api/chat", req, func(data []byte) error {
		var resp ChatResponse

		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		return fn(resp)
	})
}

func (c *Client) CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	resp := &EmbeddingResponse{}

	if err := c.do(ctx, http.MethodPost, "/api/embeddings", req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) stream(ctx context.Context, method, path string, req any, fn func([]byte) error) error {
	var body io.Reader

	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return err
		}

		body = bytes.NewReader(data)
	}

	rawurl := c.base.JoinPath(path).String()

	request, err := http.NewRequestWithContext(ctx, method, rawurl, body)
	if err != nil {
		return err
	}

	setRequestHeaders(request, "application/x-ndjson")

	resp, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	// increase the buffer size to avoid running out of space
	buf := make([]byte, 0, maxBufferSize)

	scanner.Buffer(buf, maxBufferSize)

	for scanner.Scan() {
		var errorResponse struct {
			Error string `json:"error,omitempty"`
		}

		data := scanner.Bytes()

		if err := json.Unmarshal(data, &errorResponse); err != nil {
			return err
		}

		if errorResponse.Error != "" {
			return fmt.Errorf("%s", errorResponse.Error)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return StatusError{
				StatusCode:   resp.StatusCode,
				Status:       resp.Status,
				ErrorMessage: errorResponse.Error,
			}
		}

		if err := fn(data); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) do(ctx context.Context, method, path string, req, resp any) error {
	var body io.Reader

	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return err
		}

		body = bytes.NewReader(data)
	}

	rawurl := c.base.JoinPath(path).String()

	request, err := http.NewRequestWithContext(ctx, method, rawurl, body)
	if err != nil {
		return err
	}

	setRequestHeaders(request, "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := checkError(response, data); err != nil {
		return err
	}

	if len(data) > 0 && resp != nil {
		if err := json.Unmarshal(data, resp); err != nil {
			return err
		}
	}

	return nil
}

func checkError(resp *http.Response, data []byte) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	apiError := StatusError{StatusCode: resp.StatusCode}

	if err := json.Unmarshal(data, &apiError); err != nil {
		// Use the full body as the message if we fail to decode a response.
		apiError.ErrorMessage = string(data)
	}

	return apiError
}

func setRequestHeaders(r *http.Request, accept string) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", accept)
	r.Header.Set("User-Agent", fmt.Sprintf("lc/ (%s %s) Go/%s",
		runtime.GOARCH, runtime.GOOS, runtime.Version(),
	))
}
