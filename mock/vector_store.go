package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.VectorStore = VectorStore{}

type VectorStore struct {
	AddDocumentsFunc     func(ctx context.Context, docs []llm.Document, options ...llm.VectorStoreOption) ([]string, error)
	SimilaritySearchFunc func(ctx context.Context, query string, numDocs int, options ...llm.VectorStoreOption) ([]llm.Document, error)
}

func (vs VectorStore) AddDocuments(ctx context.Context, docs []llm.Document, options ...llm.VectorStoreOption) ([]string, error) {
	return vs.AddDocumentsFunc(ctx, docs, options...)
}

func (vs VectorStore) SimilaritySearch(ctx context.Context, query string, numDocs int, options ...llm.VectorStoreOption) ([]llm.Document, error) {
	return vs.SimilaritySearchFunc(ctx, query, numDocs, options...)
}
