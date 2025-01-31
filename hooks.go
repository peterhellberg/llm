package llm

import "context"

// Hooks is the interface that allows for
// hooking into specific parts of an LLM application.
type Hooks interface {
	AgentHooks
	ChainHooks
	ProviderHooks
	RetrieverHooks
	ToolHooks
}

type AgentHooker interface {
	AgentHooks() AgentHooks
}

type ChainHooker interface {
	ChainHooks() ChainHooks
}

type AgentHooks interface {
	AgentAction(ctx context.Context, action AgentAction)
	AgentFinish(ctx context.Context, finish AgentFinish)
	StreamingFunc(ctx context.Context, chunk []byte)
}

type ChainHooks interface {
	ChainStart(ctx context.Context, inputs map[string]any)
	ChainEnd(ctx context.Context, outputs map[string]any)
	ChainError(ctx context.Context, err error)
	StreamingFunc(ctx context.Context, chunk []byte)
}

type ProviderHooks interface {
	ProviderStart(ctx context.Context, prompts []string)
	ProviderGenerateContentStart(ctx context.Context, ms []Message)
	ProviderGenerateContentEnd(ctx context.Context, res *ContentResponse)
	ProviderError(ctx context.Context, err error)
	StreamingFunc(ctx context.Context, chunk []byte)
}

type RetrieverHooks interface {
	RetrieverStart(ctx context.Context, query string)
	RetrieverEnd(ctx context.Context, query string, documents []Document)
}

type ToolHooks interface {
	ToolStart(ctx context.Context, input string)
	ToolEnd(ctx context.Context, output string)
	ToolError(ctx context.Context, err error)
}
