package qdrantstore

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/vectorstores/qdrantstore/internal/qdrant"
)

var _ llm.VectorStore = &Store{}

type Store struct {
	client   *qdrant.Client
	options  qdrant.Options
	embedder llm.Embedder
}

func New(options ...Option) (*Store, error) {
	s := &Store{
		options: qdrant.DefaultOptions(),
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if s.embedder == nil {
		return nil, fmt.Errorf("%w: missing embedder", llm.ErrInvalidOptions)
	}

	var err error

	s.client, err = qdrant.NewClient(s.options)

	return s, err
}

func (s *Store) AddDocuments(ctx context.Context, docs []llm.Document, _ ...llm.VectorStoreOption) ([]string, error) {
	texts := make([]string, 0, len(docs))

	for _, doc := range docs {
		texts = append(texts, doc.PageContent)
	}

	vectors, err := s.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}

	if len(vectors) != len(docs) {
		return nil, fmt.Errorf("number of vectors from embedder does not match number of documents")
	}

	metadatas := make([]map[string]interface{}, 0, len(docs))

	for i := 0; i < len(docs); i++ {
		metadata := make(map[string]interface{}, len(docs[i].Metadata))

		for key, value := range docs[i].Metadata {
			metadata[key] = value
		}

		metadata[s.options.ContentKey] = texts[i]

		metadatas = append(metadatas, metadata)
	}

	return s.client.UpsertPoints(ctx, vectors, metadatas)
}

func (s *Store) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...llm.VectorStoreOption) ([]llm.Document, error) {
	opts := s.getOptions(options...)

	filters := s.getFilters(opts)

	scoreThreshold, err := s.getScoreThreshold(opts)
	if err != nil {
		return nil, err
	}

	vector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return s.client.SearchPoints(ctx, vector, numDocuments, scoreThreshold, filters)
}

func (s *Store) ClearCollection(ctx context.Context) error {
	return s.client.ClearCollection(ctx)
}

func (s *Store) getScoreThreshold(opts llm.VectorStoreOptions) (float32, error) {
	if opts.ScoreThreshold < 0 || opts.ScoreThreshold > 1 {
		return 0, fmt.Errorf("score threshold must be between 0 and 1")
	}

	return opts.ScoreThreshold, nil
}

func (s *Store) getFilters(opts llm.VectorStoreOptions) any {
	if opts.Filters != nil {
		return opts.Filters
	}

	return nil
}

func (s *Store) getOptions(options ...llm.VectorStoreOption) llm.VectorStoreOptions {
	opts := llm.VectorStoreOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}
