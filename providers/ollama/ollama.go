package ollama

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/ollama/internal/ollama"
)

var (
	_ llm.Provider       = (*Provider)(nil)
	_ llm.EmbedderClient = (*Provider)(nil)
)

var (
	ErrEmptyResponse                    = fmt.Errorf("no response")
	ErrIncompleteEmbedding              = fmt.Errorf("not all input got embedded")
	ErrExpectingSingleText              = fmt.Errorf("expecting a single Text content")
	ErrOnlySupportsTextAndBinaryContent = fmt.Errorf("only supports Text and BinaryContent parts right now")
)

// Provider is an llm.Provider implementation for Ollama.
type Provider struct {
	client *ollama.Client
	Options
}

// New creates a new ollama llm.Provider implementation.
func New(options ...Option) (*Provider, error) {
	o := Options{}

	for _, option := range options {
		if err := option(&o); err != nil {
			return nil, err
		}
	}

	client, err := ollama.NewClient(o.ollamaServerURL, o.httpClient)
	if err != nil {
		return nil, err
	}

	return &Provider{
		client:  client,
		Options: o,
	}, nil
}

// Call Implement the call interface for Provider.
func (p *Provider) Call(ctx context.Context, prompt string, options ...llm.ContentOption) (string, error) {
	return llm.Call(ctx, p, prompt, options...)
}

// GenerateContent implements the Model interface.
func (p *Provider) GenerateContent(ctx context.Context, messages []llm.Message, options ...llm.ContentOption) (*llm.ContentResponse, error) {
	if p.hooks != nil {
		p.hooks.ProviderGenerateContentStart(ctx, messages)
	}

	opts := llm.ContentOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	// Override LLM model if set as llms.CallOption
	model := p.model

	if opts.Model != "" {
		model = opts.Model
	}

	format := p.format

	if opts.JSONMode {
		format = "json"
	}

	// Get our ollamaOptions from llm.CallOptions
	ollamaOptions := makeOllamaOptions(p.ollamaOptions, opts)

	// Our input is a slice of llm.Message, each of which potentially has
	// a slice of Part that could be text, images etc.
	// We have to convert it to a format Ollama undestands: ChatRequest, which
	// has a sequence of Message, each of which has a role and content - single
	// text + potential images.
	ollamaMessages, err := makeOllamaMessages(messages)
	if err != nil {
		return nil, err
	}

	req := &ollama.ChatRequest{
		Model:    model,
		Format:   format,
		Messages: ollamaMessages,
		Options:  ollamaOptions,
		Stream:   opts.StreamingFunc != nil,
	}

	keepAlive := p.keepAlive

	if keepAlive != "" {
		req.KeepAlive = keepAlive
	}

	var fn ollama.ChatResponseFunc

	streamedResponse := ""

	var resp ollama.ChatResponse

	fn = func(response ollama.ChatResponse) error {
		if opts.StreamingFunc != nil && response.Message != nil {
			if err := opts.StreamingFunc(ctx, []byte(response.Message.Content)); err != nil {
				return err
			}
		}

		if response.Message != nil {
			streamedResponse += response.Message.Content
		}

		if !req.Stream || response.Done {
			resp = response
			resp.Message = &ollama.Message{
				Role:    "assistant",
				Content: streamedResponse,
			}
		}
		return nil
	}

	if err := p.client.GenerateChat(ctx, req, fn); err != nil {
		if p.hooks != nil {
			p.hooks.ProviderError(ctx, err)
		}

		return nil, err
	}

	choices := []*llm.ContentChoice{
		{
			Content: resp.Message.Content,
			GenerationInfo: map[string]any{
				"CompletionTokens": resp.EvalCount,
				"PromptTokens":     resp.PromptEvalCount,
				"TotalTokens":      resp.EvalCount + resp.PromptEvalCount,
			},
		},
	}

	response := &llm.ContentResponse{
		Choices: choices,
	}

	if p.hooks != nil {
		p.hooks.ProviderGenerateContentEnd(ctx, response)
	}

	return response, nil
}

func (p *Provider) CreateEmbedding(ctx context.Context, inputs []string) ([][]float32, error) {
	embeddings := [][]float32{}

	for _, input := range inputs {
		req := &ollama.EmbeddingRequest{
			Prompt: input,
			Model:  p.model,
		}

		if p.keepAlive != "" {
			req.KeepAlive = p.keepAlive
		}

		embedding, err := p.client.CreateEmbedding(ctx, req)
		if err != nil {
			return nil, err
		}

		if len(embedding.Embedding) == 0 {
			return nil, ErrEmptyResponse
		}

		embeddings = append(embeddings, embedding.Embedding)
	}

	if len(inputs) != len(embeddings) {
		return embeddings, ErrIncompleteEmbedding
	}

	return embeddings, nil
}

func makeOllamaMessages(llmMessages []llm.Message) ([]*ollama.Message, error) {
	ollamaMessages := make([]*ollama.Message, 0, len(llmMessages))

	for _, mc := range llmMessages {
		msg := &ollama.Message{
			Role: typeToRole(mc.Role),
		}

		// Look at all the parts in mc; expect to find a single Text part and
		// any number of binary parts.
		var text string

		foundText := false

		var images []ollama.ImageData

		for _, p := range mc.Parts {
			switch pt := p.(type) {
			case llm.TextContent:
				if foundText {
					return nil, ErrExpectingSingleText
				}

				foundText = true
				text = pt.Text
			case llm.BinaryContent:
				images = append(images, ollama.ImageData(pt.Data))
			default:
				return nil, ErrOnlySupportsTextAndBinaryContent
			}
		}

		msg.Content = text
		msg.Images = images

		ollamaMessages = append(ollamaMessages, msg)
	}

	return ollamaMessages, nil
}

func makeOllamaOptions(o ollama.Options, co llm.ContentOptions) ollama.Options {
	o.NumPredict = co.MaxTokens
	o.Temperature = float32(co.Temperature)
	o.Stop = co.StopWords
	o.TopK = co.TopK
	o.TopP = float32(co.TopP)
	o.Seed = co.Seed
	o.RepeatPenalty = float32(co.RepetitionPenalty)
	o.FrequencyPenalty = float32(co.FrequencyPenalty)
	o.PresencePenalty = float32(co.PresencePenalty)

	return o
}

func typeToRole(cmt llm.ChatMessageType) string {
	switch cmt {
	case llm.ChatMessageTypeSystem:
		return "system"
	case llm.ChatMessageTypeAI:
		return "assistant"
	case llm.ChatMessageTypeHuman:
		fallthrough
	case llm.ChatMessageTypeGeneric:
		return "user"
	case llm.ChatMessageTypeFunction:
		return "function"
	case llm.ChatMessageTypeTool:
		return "tool"
	}
	return ""
}
