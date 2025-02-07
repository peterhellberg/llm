package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/ollama"
	"github.com/peterhellberg/llm/vectorstores/qdrantstore"
)

func main() {
	ctx := context.Background()
	env := llm.NewEnv(os.Getenv)

	if err := run(ctx, env); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, env llm.Env) error {
	embedderClient, err := ollama.New(
		ollama.WithHost(env.String("OLLAMA_HOST", "")),
		ollama.WithModel(env.String("OLLAMA_MODEL", "snowflake-arctic-embed:22m")),
	)
	if err != nil {
		return err
	}

	embedder, err := llm.NewEmbedder(embedderClient)
	if err != nil {
		return err
	}

	// Create a new Qdrant vector store.
	store, err := qdrantstore.New(
		qdrantstore.WithRawURL(env.String("QDRANT_URL", "http://localhost:6333")),
		qdrantstore.WithAPIKey(env.String("QDRANT_API_KEY", "")),
		qdrantstore.WithCollectionName(env.String("QDRANT_COLLECTION_NAME", "example")),
		qdrantstore.WithCollectionVectorSize(env.Int("QDRANT_COLLECTION_VECTOR_SIZE", 384)),
		qdrantstore.WithEmbedder(embedder),
	)
	if err != nil {
		return err
	}

	// Clear the Qdrant collection
	if err := store.ClearCollection(ctx); err != nil {
		return err
	}

	// Add documents to the Qdrant vector store.
	_, err = store.AddDocuments(ctx, places)
	if err != nil {
		return err
	}

	{ // Search for similar documents.
		docs, err := store.SimilaritySearch(ctx, "england", 1)
		if err != nil {
			return err
		}

		fmt.Println("\nengland:\n", docs)
	}

	{ // Search for similar documents using score threshold.
		docs, err := store.SimilaritySearch(ctx, "american places", 3,
			llm.VectorStoreWithScoreThreshold(0.80),
		)
		if err != nil {
			return err
		}

		fmt.Println("\namerican places:\n", docs)
	}

	{ // Search for similar documents using score threshold and metadata filter.
		filter := map[string]any{
			"must": []map[string]any{
				{
					"key": "area",
					"range": map[string]any{
						"lte": 2000,
					},
				},
			},
		}

		docs, err := store.SimilaritySearch(ctx, "only cities in south america", 10,
			llm.VectorStoreWithScoreThreshold(0.80),
			llm.VectorStoreWithFilters(filter),
		)
		if err != nil {
			return err
		}

		fmt.Println("\nonly cities in south america:\n", docs)
	}

	return nil
}

var places = []llm.Document{
	{
		PageContent: "A city in texas",
		Metadata: map[string]any{
			"area": 3251,
		},
	},
	{
		PageContent: "A country in Asia",
		Metadata: map[string]any{
			"area": 2342,
		},
	},
	{
		PageContent: "A country in South America",
		Metadata: map[string]any{
			"area": 432,
		},
	},
	{
		PageContent: "An island nation in the Pacific Ocean",
		Metadata: map[string]any{
			"area": 6531,
		},
	},
	{
		PageContent: "A mountainous country in Europe",
		Metadata: map[string]any{
			"area": 2211,
		},
	},
	{
		PageContent: "A lost city in the Amazon",
		Metadata: map[string]any{
			"area": 1223,
		},
	},
	{
		PageContent: "A city in England",
		Metadata: map[string]any{
			"area": 4324,
		},
	},
}
