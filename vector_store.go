package llm

import "context"

// VectorStore is the interface for saving and querying documents in the form of vector embeddings.
type VectorStore interface {
	AddDocuments(ctx context.Context, docs []Document, options ...VectorStoreOption) ([]string, error)
	SimilaritySearch(ctx context.Context, query string, numDocs int, options ...VectorStoreOption) ([]Document, error)
}

// VectorStoreOption is a function that configures a VectorStoreOptions.
type VectorStoreOption func(*VectorStoreOptions)

// VectorStoreOptions is a set of options for similarity search and add documents.
type VectorStoreOptions struct {
	NameSpace      string
	ScoreThreshold float32
	Filters        any
	Embedder       Embedder
	Deduplicater   func(context.Context, Document) bool
}

// VectorStoreWithNameSpace returns a VectorStoreOption for setting the name space.
func VectorStoreWithNameSpace(nameSpace string) VectorStoreOption {
	return func(o *VectorStoreOptions) {
		o.NameSpace = nameSpace
	}
}

func VectorStoreWithScoreThreshold(scoreThreshold float32) VectorStoreOption {
	return func(o *VectorStoreOptions) {
		o.ScoreThreshold = scoreThreshold
	}
}

// VectorStoreWithFilters searches can be limited based on metadata filters. Searches with  metadata
// filters retrieve exactly the number of nearest-neighbors results that match the filters. In
// most cases the search latency will be lower than unfiltered searches
// See https://docs.pinecone.io/docs/metadata-filtering
func VectorStoreWithFilters(filters any) VectorStoreOption {
	return func(o *VectorStoreOptions) {
		o.Filters = filters
	}
}

// VectorStoreWithEmbedder returns a VectorStoreOption for setting the embedder that could be used when
// adding documents or doing similarity search (instead the embedder from the Store context)
// this is useful when we are using multiple LLMs with single vectorstore.
func VectorStoreWithEmbedder(embedder Embedder) VectorStoreOption {
	return func(o *VectorStoreOptions) {
		o.Embedder = embedder
	}
}

// VectorStoreWithDeduplicater returns a VectorStoreOption for setting the deduplicater that could be used
// when adding documents. This is useful to prevent wasting time on creating an embedding
// when one already exists.
func VectorStoreWithDeduplicater(fn func(ctx context.Context, doc Document) bool) VectorStoreOption {
	return func(o *VectorStoreOptions) {
		o.Deduplicater = fn
	}
}
