package textloader

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Loader = Loader{}

// Loader loads text data from an io.Reader.
type Loader struct {
	r io.Reader
}

// New creates a new text loader with an io.Reader.
func New(r io.Reader) Loader {
	return Loader{
		r: r,
	}
}

// Load reads from the io.Reader and returns a single document with the data.
func (l Loader) Load(_ context.Context) ([]llm.Document, error) {
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, l.r); err != nil {
		return nil, err
	}

	return []llm.Document{
		{
			PageContent: strings.TrimSuffix(buf.String(), "\n"),
			Metadata:    map[string]any{},
		},
	}, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple documents using a text splitter.
func (l Loader) LoadAndSplit(ctx context.Context, splitter llm.TextSplitter) ([]llm.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return llm.SplitDocuments(splitter, docs)
}
