package llm

import "context"

// Embedder is the interface for creating vector embeddings from texts.
type Embedder interface {
	// EmbedDocuments returns a vector for each text.
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
	// EmbedQuery embeds a single text.
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
}

// EmbedderClient is the interface LLM clients implement for embeddings.
type EmbedderClient interface {
	CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error)
}

// EmbedderClientFunc is an adapter to allow the use of ordinary functions as Embedder Clients. If
// `f` is a function with the appropriate signature, `EmbedderClientFunc(f)` is an `EmbedderClient`
// that calls `f`.
type EmbedderClientFunc func(ctx context.Context, texts []string) ([][]float32, error)

func (e EmbedderClientFunc) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	return e(ctx, texts)
}
