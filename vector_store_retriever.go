package llm

import "context"

// VectorStoreRetriever is a retriever for vector stores.
type VectorStoreRetriever struct {
	Hooks   RetrieverHooks
	vs      VectorStore
	numDocs int
	options []VectorStoreOption
}

var _ Retriever = VectorStoreRetriever{}

// RelevantDocuments returns documents using the vector store.
func (r VectorStoreRetriever) RelevantDocuments(ctx context.Context, query string) ([]Document, error) {
	if r.Hooks != nil {
		r.Hooks.RetrieverStart(ctx, query)
	}

	docs, err := r.vs.SimilaritySearch(ctx, query, r.numDocs, r.options...)
	if err != nil {
		return nil, err
	}

	if r.Hooks != nil {
		r.Hooks.RetrieverEnd(ctx, query, docs)
	}

	return docs, nil
}

// NewVectorStoreRetriever takes a vector store and returns a retriever using the vector store to retrieve documents.
func NewVectorStoreRetriever(vs VectorStore, numDocs int, options ...VectorStoreOption) VectorStoreRetriever {
	return VectorStoreRetriever{
		vs:      vs,
		numDocs: numDocs,
		options: options,
	}
}
