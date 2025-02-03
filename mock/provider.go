package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Provider = Provider{}

type Provider struct {
	GenerateContentFunc func(context.Context, []llm.Message, ...llm.ContentOption) (*llm.ContentResponse, error)
}

func (p Provider) GenerateContent(ctx context.Context, messages []llm.Message, options ...llm.ContentOption) (*llm.ContentResponse, error) {
	return p.GenerateContentFunc(ctx, messages, options...)
}
