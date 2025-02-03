package retrievalqachain

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
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

// New creates a new Chain from a retriever and a chain for combining documents.
// The chain for combining documents is expected to have the expected input values
// for the "question" and "input_documents" key.
func New(combineDocumentsChain llm.Chain, retriever llm.Retriever, options ...func(*Chain)) Chain {
	c := Chain{
		Retriever:             retriever,
		CombineDocumentsChain: combineDocumentsChain,
		InputKey:              defaultInputKey,
		ReturnSourceDocuments: false,
	}

	for _, opt := range options {
		opt(&c)
	}

	return c
}

// Call gets relevant documents from the retriever and gives them to the combine documents chain.
func (c Chain) Call(ctx context.Context, values map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	query, ok := values[c.InputKey].(string)
	if !ok {
		return nil, fmt.Errorf("%w: %w", llm.ErrInvalidInputValues, llm.ErrInputValuesWrongType)
	}

	docs, err := c.Retriever.RelevantDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	inputValues := map[string]any{
		"question":        query,
		"input_documents": docs,
	}

	result, err := llm.ChainCall(ctx, c.CombineDocumentsChain, inputValues, options...)
	if err != nil {
		return nil, err
	}

	if c.ReturnSourceDocuments {
		result[defaultSourceDocumentKey] = docs
	}

	return result, nil
}

// Memory returns empty memory.
func (c Chain) Memory() llm.Memory {
	return llm.EmptyMemory{}
}

func (c Chain) InputKeys() []string {
	return []string{c.InputKey}
}

func (c Chain) OutputKeys() []string {
	outputKeys := append([]string{},
		c.CombineDocumentsChain.OutputKeys()...,
	)

	if c.ReturnSourceDocuments {
		outputKeys = append(outputKeys, defaultSourceDocumentKey)
	}

	return outputKeys
}
