package csvloader

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Loader = Loader{}

// Loader implements the llm.Loader interface for CSV.
type Loader struct {
	r       io.Reader
	columns []string
}

// New creates a new CSV loader with an io.Reader and optional column names for filtering.
func New(r io.Reader, columns ...string) Loader {
	return Loader{
		r:       r,
		columns: columns,
	}
}

// Load reads from the io.Reader and returns a single document with the data.
func (c Loader) Load(_ context.Context) ([]llm.Document, error) {
	var (
		header []string
		docs   []llm.Document
		rown   int
	)

	rd := csv.NewReader(c.r)

	for {
		record, err := rd.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(header) == 0 {
			header = append(header, record...)
			continue
		}

		var content []string

		for i, value := range record {
			if len(c.columns) > 0 && !slices.Contains(c.columns, header[i]) {
				continue
			}

			line := fmt.Sprintf("%s: %s", header[i], value)

			content = append(content, line)
		}

		rown++

		docs = append(docs, llm.Document{
			PageContent: strings.Join(content, "\n"),
			Metadata: map[string]any{
				"row": rown,
			},
		})
	}

	return docs, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple documents using a text splitter.
func (c Loader) LoadAndSplit(ctx context.Context, splitter llm.TextSplitter) ([]llm.Document, error) {
	docs, err := c.Load(ctx)
	if err != nil {
		return nil, err
	}

	return llm.SplitDocuments(splitter, docs)
}
