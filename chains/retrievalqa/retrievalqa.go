package retrievalqa

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/chains"
	"github.com/peterhellberg/llm/memory"
)

var _ llm.Chain = Chain{}

const (
	defaultInputKey          = "query"
	defaultSourceDocumentKey = "source_documents"
)

// Chain used for question-answering against a retriever.
// First the chain gets documents from the retriever, then the documents
// and the query is used as input to another chain. Typically, that chain
// combines the documents into a prompt that is sent to an LLM.

type Chain struct {
	// Retriever used to retrieve the relevant documents.
	Retriever llm.Retriever

	// The chain the documents and query is given to.
	CombineDocumentsChain llm.Chain

	// The input key to get the query from, by default "query".
	InputKey string

	// If the chain should return the documents used by the combine
	// documents chain in the "source_documents" key.
	ReturnSourceDocuments bool
}

// New creates a new Chain from a retriever and a chain for
// combining documents. The chain for combining documents is expected to
// have the expected input values for the "question" and "input_documents"
// key.
func New(combineDocumentsChain llm.Chain, retriever llm.Retriever) Chain {
	return Chain{
		Retriever:             retriever,
		CombineDocumentsChain: combineDocumentsChain,
		InputKey:              defaultInputKey,
		ReturnSourceDocuments: false,
	}
}

// Call gets relevant documents from the retriever and gives them to the combine
// documents chain.
func (c Chain) Call(ctx context.Context, values map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	query, ok := values[c.InputKey].(string)
	if !ok {
		return nil, fmt.Errorf("%w: %w", llm.ErrInvalidInputValues, llm.ErrInputValuesWrongType)
	}

	docs, err := c.Retriever.RelevantDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := chains.Call(ctx, c.CombineDocumentsChain, map[string]any{
		"question":        query,
		"input_documents": docs,
	}, options...)
	if err != nil {
		return nil, err
	}

	if c.ReturnSourceDocuments {
		result[defaultSourceDocumentKey] = docs
	}

	return result, nil
}

func (c Chain) GetMemory() llm.Memory {
	return memory.Empty{}
}

func (c Chain) GetInputKeys() []string {
	return []string{c.InputKey}
}

func (c Chain) GetOutputKeys() []string {
	outputKeys := append([]string{}, c.CombineDocumentsChain.GetOutputKeys()...)

	if c.ReturnSourceDocuments {
		outputKeys = append(outputKeys, defaultSourceDocumentKey)
	}

	return outputKeys
}
