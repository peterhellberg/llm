package openai

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/openai/internal/openai"
)

var (
	_ llm.Provider       = (*Provider)(nil)
	_ llm.EmbedderClient = (*Provider)(nil)
)

// Provider is an llm.Provider implementation for OpenAI.
type Provider struct {
	client *openai.Client
	hooks  llm.ProviderHooks
}

// New creates a new OpenAI llm.Provider implementation.
func New(opts ...Option) (*Provider, error) {
	opt, c, err := newClient(os.Getenv, opts...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		client: c,
		hooks:  opt.hooks,
	}, err
}

// Call requests a completion for the given prompt.
func (o *Provider) Call(ctx context.Context, prompt string, options ...llm.ContentOption) (string, error) {
	return llm.Call(ctx, o, prompt, options...)
}

// GenerateContent implements the llm.Provider interface.
func (o *Provider) GenerateContent(ctx context.Context, messages []llm.Message, options ...llm.ContentOption) (*llm.ContentResponse, error) {
	if o.hooks != nil {
		o.hooks.ProviderGenerateContentStart(ctx, messages)
	}

	opts := llm.ContentOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	chatMsgs := make([]*openai.ChatMessage, 0, len(messages))

	for _, mc := range messages {
		msg := &openai.ChatMessage{
			MultiContent: mc.Parts,
		}

		switch mc.Role {
		case llm.ChatMessageTypeSystem:
			msg.Role = roleSystem
		case llm.ChatMessageTypeAI:
			msg.Role = roleAssistant
		case llm.ChatMessageTypeHuman:
			msg.Role = roleUser
		case llm.ChatMessageTypeGeneric:
			msg.Role = roleUser
		case llm.ChatMessageTypeFunction:
			msg.Role = roleFunction
		case llm.ChatMessageTypeTool:
			msg.Role = roleTool

			if len(mc.Parts) != 1 {
				return nil, fmt.Errorf("expected exactly one part for role %v, got %v", mc.Role, len(mc.Parts))
			}

			switch p := mc.Parts[0].(type) {
			case llm.ToolCallResponse:
				msg.ToolCallID = p.ToolCallID
				msg.Content = p.Content
			default:
				return nil, fmt.Errorf("expected part of type ToolCallResponse for role %v, got %T", mc.Role, mc.Parts[0])
			}
		default:
			return nil, fmt.Errorf("role %v not supported", mc.Role)
		}

		newParts, toolCalls := ExtractToolParts(msg)
		msg.MultiContent = newParts
		msg.ToolCalls = toolCallsFromToolCalls(toolCalls)

		chatMsgs = append(chatMsgs, msg)
	}

	req := &openai.ChatRequest{
		Model:            opts.Model,
		StopWords:        opts.StopWords,
		Messages:         chatMsgs,
		StreamingFunc:    opts.StreamingFunc,
		Temperature:      opts.Temperature,
		N:                opts.N,
		FrequencyPenalty: opts.FrequencyPenalty,
		PresencePenalty:  opts.PresencePenalty,

		MaxCompletionTokens: opts.MaxTokens,

		ToolChoice: opts.ToolChoice,
		Seed:       opts.Seed,
		Metadata:   opts.Metadata,
	}
	if opts.JSONMode {
		req.ResponseFormat = ResponseFormatJSON
	}

	// if opts.Tools is not empty, append them to req.Tools
	for _, tool := range opts.Tools {
		t, err := toolFromTool(tool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert llm tool to openai tool: %w", err)
		}

		req.Tools = append(req.Tools, t)
	}

	// if o.client.ResponseFormat is set, use it for the request
	if o.client.ResponseFormat != nil {
		req.ResponseFormat = o.client.ResponseFormat
	}

	result, err := o.client.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	choices := make([]*llm.ContentChoice, len(result.Choices))

	for i, c := range result.Choices {
		choices[i] = &llm.ContentChoice{
			Content:    c.Message.Content,
			StopReason: fmt.Sprint(c.FinishReason),
			GenerationInfo: map[string]any{
				"CompletionTokens": result.Usage.CompletionTokens,
				"PromptTokens":     result.Usage.PromptTokens,
				"TotalTokens":      result.Usage.TotalTokens,
				"ReasoningTokens":  result.Usage.CompletionTokensDetails.ReasoningTokens,
			},
		}

		for _, tool := range c.Message.ToolCalls {
			choices[i].ToolCalls = append(choices[i].ToolCalls, llm.ToolCall{
				ID:   tool.ID,
				Type: string(tool.Type),
				FunctionCall: &llm.FunctionCall{
					Name:      tool.Function.Name,
					Arguments: tool.Function.Arguments,
				},
			})
		}

		// populate legacy single-function call field for backwards compatibility
		if len(choices[i].ToolCalls) > 0 {
			choices[i].FuncCall = choices[i].ToolCalls[0].FunctionCall
		}
	}

	response := &llm.ContentResponse{Choices: choices}

	if o.hooks != nil {
		o.hooks.ProviderGenerateContentEnd(ctx, response)
	}
	return response, nil
}

// CreateEmbedding creates embeddings for the given input texts.
func (o *Provider) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
	embeddings, err := o.client.CreateEmbedding(ctx, &openai.EmbeddingRequest{
		Input: inputTexts,
		Model: o.client.EmbeddingModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create openai embeddings: %w", err)
	}

	if len(embeddings) == 0 {
		return nil, ErrEmptyResponse
	}

	if len(inputTexts) != len(embeddings) {
		return embeddings, ErrUnexpectedResponseLength
	}

	return embeddings, nil
}

// ExtractToolParts extracts the tool parts from a message.
func ExtractToolParts(msg *openai.ChatMessage) ([]llm.ContentPart, []llm.ToolCall) {
	var (
		content   []llm.ContentPart
		toolCalls []llm.ToolCall
	)

	for _, part := range msg.MultiContent {
		switch p := part.(type) {
		case llm.TextContent:
			content = append(content, p)
		case llm.ImageURLContent:
			content = append(content, p)
		case llm.BinaryContent:
			content = append(content, p)
		case llm.ToolCall:
			toolCalls = append(toolCalls, p)
		}
	}

	return content, toolCalls
}

// toolFromTool converts an llm.Tool to a Tool.
func toolFromTool(t llm.Tool) (openai.Tool, error) {
	tool := openai.Tool{
		Type: openai.ToolType(t.Type),
	}

	switch t.Type {
	case string(openai.ToolTypeFunction):
		tool.Function = openai.FunctionDefinition{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
			Strict:      t.Function.Strict,
		}
	default:
		return openai.Tool{}, fmt.Errorf("tool type %v not supported", t.Type)
	}

	return tool, nil
}

// toolCallsFromToolCalls converts a slice of llm.ToolCall to a slice of ToolCall.
func toolCallsFromToolCalls(tcs []llm.ToolCall) []openai.ToolCall {
	toolCalls := make([]openai.ToolCall, len(tcs))

	for i, tc := range tcs {
		toolCalls[i] = toolCallFromToolCall(tc)
	}

	return toolCalls
}

// toolCallFromToolCall converts an llm.ToolCall to a ToolCall.
func toolCallFromToolCall(tc llm.ToolCall) openai.ToolCall {
	return openai.ToolCall{
		ID:   tc.ID,
		Type: openai.ToolType(tc.Type),
		Function: openai.ToolFunction{
			Name:      tc.FunctionCall.Name,
			Arguments: tc.FunctionCall.Arguments,
		},
	}
}

// newClient creates an instance of the internal client.
func newClient(getenv llm.Getenv, opts ...Option) (*options, *openai.Client, error) {
	options := &options{
		token:        getenv(tokenEnvVarName),
		model:        getenv(modelEnvVarName),
		baseURL:      getEnvs(getenv, baseURLEnvVarName, baseAPIBaseEnvVarName),
		organization: getenv(organizationEnvVarName),
		apiType:      APIType(openai.APITypeOpenAI),
		httpClient:   http.DefaultClient,
	}

	for _, opt := range opts {
		opt(options)
	}

	if openai.IsAzure(openai.APIType(options.apiType)) && options.apiVersion == "" {
		options.apiVersion = DefaultAPIVersion

		if options.model == "" {
			return options, nil, ErrMissingAzureModel
		}

		if options.embeddingModel == "" {
			return options, nil, ErrMissingAzureEmbeddingModel
		}
	}

	if len(options.token) == 0 {
		return options, nil, ErrMissingToken
	}

	cli, err := openai.New(
		options.token,
		options.model,
		options.baseURL,
		options.organization,
		openai.APIType(options.apiType),
		options.apiVersion,
		options.httpClient,
		options.embeddingModel,
		options.responseFormat,
	)

	return options, cli, err
}

func getEnvs(getenv func(string) string, keys ...string) string {
	for _, key := range keys {
		if val := getenv(key); val != "" {
			return val
		}
	}

	return ""
}
