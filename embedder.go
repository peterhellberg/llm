package llm

import (
	"context"
	"fmt"
	"strings"
)

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

// NewEmbedder creates a new Embedder from the given EmbedderClient, with
// some options that affect how embedding will be done.
func NewEmbedder(client EmbedderClient, opts ...EmbedderOption) (Embedder, error) {
	const (
		defaultBatchSize     = 512
		defaultStripNewLines = true
	)

	e := &embedder{
		client:        client,
		stripNewLines: defaultStripNewLines,
		batchSize:     defaultBatchSize,
	}

	for _, opt := range opts {
		opt(e)
	}
	return e, nil
}

type EmbedderOption func(p *embedder)

// WithStripNewLines is an option for specifying the should it strip new lines.
func WithStripNewLines(stripNewLines bool) EmbedderOption {
	return func(p *embedder) {
		p.stripNewLines = stripNewLines
	}
}

// WithBatchSize is an option for specifying the batch size.
func WithBatchSize(batchSize int) EmbedderOption {
	return func(p *embedder) {
		p.batchSize = batchSize
	}
}

type embedder struct {
	client EmbedderClient

	stripNewLines bool
	batchSize     int
}

// EmbedQuery embeds a single text.
func (ei *embedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	if ei.stripNewLines {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	emb, err := ei.client.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("error embedding query: %w", err)
	}

	return emb[0], nil
}

// EmbedDocuments creates one vector embedding for each of the texts.
func (ei *embedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	texts = maybeRemoveNewLines(texts, ei.stripNewLines)
	return batchedEmbed(ctx, ei.client, texts, ei.batchSize)
}

func maybeRemoveNewLines(texts []string, removeNewLines bool) []string {
	if !removeNewLines {
		return texts
	}

	for i := 0; i < len(texts); i++ {
		texts[i] = strings.ReplaceAll(texts[i], "\n", " ")
	}

	return texts
}

// batchedEmbed creates embeddings for the given input texts, batching them
// into batches of batchSize if needed.
func batchedEmbed(ctx context.Context, embedder EmbedderClient, texts []string, batchSize int) ([][]float32, error) {
	batchedTexts := batchTexts(texts, batchSize)

	emb := make([][]float32, 0, len(texts))
	for _, batch := range batchedTexts {
		curBatchEmbeddings, err := embedder.CreateEmbedding(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("error embedding batch: %w", err)
		}
		emb = append(emb, curBatchEmbeddings...)
	}

	return emb, nil
}

// batchTexts splits strings by the length batchSize.
func batchTexts(texts []string, batchSize int) [][]string {
	batchedTexts := make([][]string, 0, len(texts)/batchSize+1)

	for i := 0; i < len(texts); i += batchSize {
		batchedTexts = append(batchedTexts, texts[i:minInt([]int{i + batchSize, len(texts)})])
	}

	return batchedTexts
}

// minInt returns the minimum value in nums.
// If nums is empty, it returns 0.
func minInt(nums []int) int {
	var m int
	for idx := 0; idx < len(nums); idx++ {
		item := nums[idx]
		if idx == 0 {
			m = item
			continue
		}
		if item < m {
			m = item
		}
	}
	return m
}
