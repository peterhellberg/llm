package llm

import "context"

// VectorStore is the interface for saving and querying documents in the form of vector embeddings.
type VectorStore interface {
	AddDocuments(ctx context.Context, docs []Document, options ...VectorStoreOption) ([]string, error)
	SimilaritySearch(ctx context.Context, query string, numDocs int, options ...VectorStoreOption) ([]Document, error)
}
