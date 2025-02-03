package llm

import "context"

var _ Retriever = VectorStoreRetriever{}

// VectorStoreRetriever is a retriever for vector stores.
type VectorStoreRetriever struct {
	Hooks   RetrieverHooks
	vs      VectorStore
	numDocs int
	options []VectorStoreOption
}

// NewVectorStoreRetriever takes a vector store and returns a retriever using the vector store to retrieve documents.
func NewVectorStoreRetriever(vs VectorStore, numDocs int, options ...VectorStoreRetrieverOption) VectorStoreRetriever {
	vsr := VectorStoreRetriever{
		vs:      vs,
		numDocs: numDocs,
	}

	for _, option := range options {
		option(&vsr)
	}

	return vsr
}

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

// VectorStoreRetrieverOption is a function that configures a VectorStoreRetriever.
type VectorStoreRetrieverOption func(*VectorStoreRetriever)

func VectorStoreRetrieverWithHooks(hooks RetrieverHooks) VectorStoreRetrieverOption {
	return func(o *VectorStoreRetriever) {
		o.Hooks = hooks
	}
}

func VectorStoreRetrieverWithVectorStoreOptions(options ...VectorStoreOption) VectorStoreRetrieverOption {
	return func(o *VectorStoreRetriever) {
		o.options = options
	}
}
