package textloader

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	_ "embed"

	"github.com/peterhellberg/llm/splitters/charsplitter"
)

//go:embed testdata/data.txt
var data []byte

func TestLoader(t *testing.T) {
	ctx := context.Background()

	t.Run("Load", func(t *testing.T) {
		l := New(bytes.NewReader(data))

		docs, err := l.Load(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(docs), 1; got != want {
			t.Fatalf("len(docs) = %d, want %d", got, want)
		}

		doc := docs[0]

		if got, want := doc.PageContent, "Foo Bar Baz"; got != want {
			t.Fatalf("doc.PageContent = %q, want %q", got, want)
		}

		if got, want := doc.Metadata, map[string]any{}; !reflect.DeepEqual(got, want) {
			t.Fatalf("doc.Metadata = %v, want %v", got, want)
		}
	})

	t.Run("LoadAndSplit", func(t *testing.T) {
		l := New(bytes.NewReader(data))

		splitter := charsplitter.New(
			charsplitter.WithChunkSize(5),
		)

		docs, err := l.LoadAndSplit(ctx, splitter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(docs), 3; got != want {
			t.Fatalf("len(docs) = %d, want %d", got, want)
		}

		doc := docs[2]

		if got, want := doc.PageContent, "Baz"; got != want {
			t.Fatalf("doc.PageContent = %q, want %q", got, want)
		}

		if got, want := doc.Metadata, map[string]any{}; !reflect.DeepEqual(got, want) {
			t.Fatalf("doc.Metadata = %v, want %v", got, want)
		}
	})
}
