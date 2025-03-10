package llm

import "context"

// Hooks is the interface that allows for hooking
// into specific parts of an LLM application.
type Hooks interface {
	AgentHooks
	ChainHooks
	ProviderHooks
	RetrieverHooks
	ToolHooks
}

// A AgentHooker can return its AgentHooks
type AgentHooker interface {
	AgentHooks() AgentHooks
}

// A ChainHooker can return its ChainHooks
type ChainHooker interface {
	ChainHooks() ChainHooks
}

// AgentHooks contains the hooks that can be used by an agent.
type AgentHooks interface {
	AgentAction(ctx context.Context, action AgentAction)
	AgentFinish(ctx context.Context, finish AgentFinish)

	StreamingFunc(ctx context.Context, chunk []byte)
}

// ChainHooks contains the hooks that can be used by a chain.
type ChainHooks interface {
	ChainStart(ctx context.Context, inputs map[string]any)
	ChainEnd(ctx context.Context, outputs map[string]any)
	ChainError(ctx context.Context, err error)

	StreamingFunc(ctx context.Context, chunk []byte)
}

// ProviderHooks contains the hooks that can be used by a provider.
type ProviderHooks interface {
	ProviderStart(ctx context.Context, prompts []string)
	ProviderGenerateContentStart(ctx context.Context, ms []Message)
	ProviderGenerateContentEnd(ctx context.Context, res *ContentResponse)
	ProviderError(ctx context.Context, err error)

	StreamingFunc(ctx context.Context, chunk []byte)
}

// RetrieverHooks contains the hooks that can be used by a retriever.
type RetrieverHooks interface {
	RetrieverStart(ctx context.Context, query string)
	RetrieverEnd(ctx context.Context, query string, documents []Document)
}

// ToolHooks contains the hooks that can be used by a tool.
type ToolHooks interface {
	ToolStart(ctx context.Context, input string)
	ToolEnd(ctx context.Context, output string)
	ToolError(ctx context.Context, err error)
}

var _ Hooks = EmptyHooks{}

// EmptyHooks hooks that does nothing. Useful for embedding.
type EmptyHooks struct{}

func (EmptyHooks) Text(context.Context, string)          {}
func (EmptyHooks) StreamingFunc(context.Context, []byte) {}

// Agent hooks

func (EmptyHooks) AgentAction(context.Context, AgentAction) {}
func (EmptyHooks) AgentFinish(context.Context, AgentFinish) {}

// Chain hooks

func (EmptyHooks) ChainStart(context.Context, map[string]any) {}
func (EmptyHooks) ChainError(context.Context, error)          {}
func (EmptyHooks) ChainEnd(context.Context, map[string]any)   {}

// Provider hooks

func (EmptyHooks) ProviderStart(context.Context, []string)                      {}
func (EmptyHooks) ProviderGenerateContentStart(context.Context, []Message)      {}
func (EmptyHooks) ProviderGenerateContentEnd(context.Context, *ContentResponse) {}
func (EmptyHooks) ProviderError(context.Context, error)                         {}

// Retriever  hooks

func (EmptyHooks) RetrieverStart(context.Context, string)           {}
func (EmptyHooks) RetrieverEnd(context.Context, string, []Document) {}

// Tool hooks

func (EmptyHooks) ToolStart(context.Context, string) {}
func (EmptyHooks) ToolError(context.Context, error)  {}
func (EmptyHooks) ToolEnd(context.Context, string)   {}
