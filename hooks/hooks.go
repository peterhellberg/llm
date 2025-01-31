package hooks

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Hooks = Empty{}

// Empty hooks that does nothing. Useful for embedding.
type Empty struct{}

func (Empty) Text(context.Context, string)          {}
func (Empty) StreamingFunc(context.Context, []byte) {}

// Agent hooks

func (Empty) AgentAction(context.Context, llm.AgentAction) {}
func (Empty) AgentFinish(context.Context, llm.AgentFinish) {}

// Chain hooks

func (Empty) ChainStart(context.Context, map[string]any) {}
func (Empty) ChainError(context.Context, error)          {}
func (Empty) ChainEnd(context.Context, map[string]any)   {}

// Provider hooks

func (Empty) ProviderStart(context.Context, []string)                          {}
func (Empty) ProviderGenerateContentStart(context.Context, []llm.Message)      {}
func (Empty) ProviderGenerateContentEnd(context.Context, *llm.ContentResponse) {}
func (Empty) ProviderError(context.Context, error)                             {}

// Retriever  hooks

func (Empty) RetrieverStart(context.Context, string)               {}
func (Empty) RetrieverEnd(context.Context, string, []llm.Document) {}

// Tool hooks

func (Empty) ToolStart(context.Context, string) {}
func (Empty) ToolError(context.Context, error)  {}
func (Empty) ToolEnd(context.Context, string)   {}
