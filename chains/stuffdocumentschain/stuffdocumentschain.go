package stuffdocumentschain

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
)

var _ llm.Chain = Chain{}

const (
	defaultInputKey             = "input_documents"
	defaultDocumentVariableName = "context"
	defaultSeparator            = "\n\n"
)

// Chain that combines documents with a separator and uses
// the stuffed documents in an llm.Chain. The input values to the llm chain
// contains all input values given to this chain, and the stuffed document as
// a string in the key specified by the "DocumentVariableName" field that is
// by default set to "context".
type Chain struct {
	// Next is the Chain called after formatting the documents.
	Next llm.Chain

	// Input key is the input key the StuffDocuments chain expects the
	//  documents to be in.
	InputKey string

	// DocumentVariableName is the variable name used in the llm_chain to put
	// the documents in.
	DocumentVariableName string

	// Separator is the string used to join the documents.
	Separator string
}

// New creates a new stuff documents chain with an LLM chain used
// after formatting the documents.
func New(next llm.Chain, options ...func(*Chain)) Chain {
	c := Chain{
		Next: next,

		InputKey:             defaultInputKey,
		DocumentVariableName: defaultDocumentVariableName,
		Separator:            defaultSeparator,
	}

	for _, opt := range options {
		opt(&c)
	}

	return c
}

// Call handles the inner logic of the StuffDocuments chain.
func (c Chain) Call(ctx context.Context, values map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	docs, ok := values[c.InputKey].([]llm.Document)
	if !ok {
		return nil, fmt.Errorf("%w: %w", llm.ErrInvalidInputValues, llm.ErrInputValuesWrongType)
	}

	inputValues := make(map[string]any)

	for key, value := range values {
		inputValues[key] = value
	}

	inputValues[c.DocumentVariableName] = c.joinDocuments(docs)

	return llm.ChainCall(ctx, c.Next, inputValues, options...)
}

// Memory returns empty memory.
func (c Chain) Memory() llm.Memory {
	return llm.EmptyMemory{}
}

// InputKeys returns the expected input keys, by default "input_documents".
func (c Chain) InputKeys() []string {
	return []string{c.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c Chain) OutputKeys() []string {
	return append([]string{}, c.Next.OutputKeys()...)
}

// joinDocuments joins the documents with the separator.
func (c Chain) joinDocuments(docs []llm.Document) string {
	var text string

	docLen := len(docs)

	for k, doc := range docs {
		text += doc.PageContent

		if k != docLen-1 {
			text += c.Separator
		}
	}

	return text
}
