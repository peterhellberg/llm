package llm_test

import (
	"context"
	"testing"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/mock"
)

func TestVectorStoreRetriever(t *testing.T) {
	t.Run("RelevantDocuments", func(t *testing.T) {
		var (
			ctx     = context.Background()
			query   = "a test query"
			numDocs = 2

			similaritySearch = func(ctx context.Context, gotQuery string, gotNumDocs int, options ...llm.VectorStoreOption) ([]llm.Document, error) {
				if gotQuery != query {
					t.Fatalf("query = %q, want %q", gotQuery, query)
				}

				if gotNumDocs != numDocs {
					t.Fatalf("numDocs = %d, want %d", gotNumDocs, numDocs)
				}

				return []llm.Document{
					{PageContent: "Foo"},
					{PageContent: "Bar"},
				}, nil
			}

			retrieverStart = func(ctx context.Context, gotQuery string) {
				if got, want := gotQuery, query; got != want {
					t.Fatalf("query = %q, want %q", got, want)
				}

				t.Logf("RetrieverStart - query: %q", query)
			}

			retrieverEnd = func(ctx context.Context, gotQuery string, docs []llm.Document) {
				if got, want := gotQuery, query; got != want {
					t.Fatalf("query = %q, want %q", got, want)
				}

				if got, want := len(docs), numDocs; got != want {
					t.Fatalf("len(docs) = %d, want %d", got, want)
				}

				t.Logf("RetrieverEnd - query: %q docs: %+v", query, docs)
			}
		)

		llm.NewVectorStoreRetriever(
			mock.VectorStore{
				SimilaritySearchFunc: similaritySearch,
			},
			numDocs,
			llm.VectorStoreRetrieverWithHooks(
				mock.Hooks{
					RetrieverStartFunc: retrieverStart,
					RetrieverEndFunc:   retrieverEnd,
				},
			),
		).RelevantDocuments(ctx, query)
	})
}
