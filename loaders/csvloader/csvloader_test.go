package csvloader

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	_ "embed"

	"github.com/peterhellberg/llm/splitters/charsplitter"
)

//go:embed testdata/data.csv
var data []byte

func TestLoader(t *testing.T) {
	ctx := context.Background()

	t.Run("Load", func(t *testing.T) {
		l := New(bytes.NewReader(data))

		docs, err := l.Load(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(docs), 20; got != want {
			t.Fatalf("len(docs) = %d, want %d", got, want)
		}

		var (
			doc             = docs[5]
			wantPageContent = "name: Sophia Wilson\nage: 22\ncity: Berlin\ncountry: Germany"
			wantMetadata    = map[string]any{"row": 6}
		)

		if got := doc.PageContent; got != wantPageContent {
			t.Fatalf("doc.PageContent = %q, want %q", got, wantPageContent)
		}

		if got := doc.Metadata; !reflect.DeepEqual(got, wantMetadata) {
			t.Fatalf("doc.Metadata = %v, want %v", got, wantMetadata)
		}
	})

	t.Run("LoadAndSplit", func(t *testing.T) {
		l := New(bytes.NewReader(data))

		splitter := charsplitter.New(
			charsplitter.WithChunkSize(50),
		)

		docs, err := l.LoadAndSplit(ctx, splitter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(docs), 40; got != want {
			t.Fatalf("len(docs) = %d, want %d", got, want)
		}

		var (
			doc             = docs[9]
			wantPageContent = "age: 37\ncity: Toronto\ncountry: Canada"
			wantMetadata    = map[string]any{"row": 5}
		)

		if got := doc.PageContent; got != wantPageContent {
			t.Fatalf("doc.PageContent = %q, want %q", got, wantPageContent)
		}

		if got := doc.Metadata; !reflect.DeepEqual(got, wantMetadata) {
			t.Fatalf("doc.Metadata = %v, want %v", got, wantMetadata)
		}
	})
}
