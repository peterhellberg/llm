package llm

import "context"

// Provider is an interface all LLM providers must implement.
type Provider interface {
	// GenerateContent asks the model to generate content from a sequence of
	// messages. It's the most general interface for multi-modal LLMs that support
	// chat-like interactions.
	GenerateContent(ctx context.Context, messages []Message, options ...ContentOption) (*ContentResponse, error)
}

// Content is a convenience function for calling an LLM provider with a single string prompt.
func Content(ctx context.Context, provider Provider, prompt string, options ...ContentOption) (*ContentResponse, error) {
	return provider.GenerateContent(ctx, []Message{
		{
			Role: ChatMessageTypeHuman,
			Parts: []ContentPart{
				TextContent{prompt},
			},
		},
	}, options...)
}

// Call is a convenience function for calling an LLM provider with a single string prompt,
// expecting a single string response. It's useful for simple, string-only interactions
// and provides a slightly more ergonomic API than the more general [Provider.GenerateContent].
func Call(ctx context.Context, provider Provider, prompt string, options ...ContentOption) (string, error) {
	resp, err := Content(ctx, provider, prompt, options...)
	if err != nil {
		return "", err
	}

	if cs := resp.Choices; len(cs) > 0 {
		return cs[0].Content, nil
	}

	return "", ErrEmptyResponseFromProvider
}
