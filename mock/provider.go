package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var (
	_ llm.Provider       = Provider{}
	_ llm.EmbedderClient = Provider{}
)

type Provider struct {
	GenerateContentFunc func(context.Context, []llm.Message, ...llm.ContentOption) (*llm.ContentResponse, error)
	CreateEmbeddingFunc func(ctx context.Context, texts []string) ([][]float32, error)
}

func (p Provider) GenerateContent(ctx context.Context, messages []llm.Message, options ...llm.ContentOption) (*llm.ContentResponse, error) {
	return p.GenerateContentFunc(ctx, messages, options...)
}

func (p Provider) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	return p.CreateEmbeddingFunc(ctx, texts)
}
