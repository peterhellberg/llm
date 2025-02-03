package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var (
	_ llm.Hooks       = Hooks{}
	_ llm.AgentHooker = Hooks{}
	_ llm.ChainHooker = Hooks{}
)

type Hooks struct {
	StreamingFuncFunc func(ctx context.Context, chunk []byte)

	AgentActionFunc func(ctx context.Context, action llm.AgentAction)
	AgentFinishFunc func(ctx context.Context, finish llm.AgentFinish)

	ChainStartFunc func(ctx context.Context, in map[string]any)
	ChainEndFunc   func(ctx context.Context, out map[string]any)
	ChainErrorFunc func(ctx context.Context, err error)

	ProviderStartFunc                func(ctx context.Context, prompts []string)
	ProviderGenerateContentStartFunc func(ctx context.Context, messages []llm.Message)
	ProviderGenerateContentEndFunc   func(ctx context.Context, res *llm.ContentResponse)
	ProviderErrorFunc                func(ctx context.Context, err error)

	RetrieverStartFunc func(ctx context.Context, query string)
	RetrieverEndFunc   func(ctx context.Context, query string, docs []llm.Document)

	ToolStartFunc func(ctx context.Context, in string)
	ToolEndFunc   func(ctx context.Context, out string)
	ToolErrorFunc func(ctx context.Context, err error)

	AgentHooksFunc func() llm.AgentHooks
	ChainHooksFunc func() llm.ChainHooks
}

func (h Hooks) StreamingFunc(ctx context.Context, chunk []byte) {
	h.StreamingFuncFunc(ctx, chunk)
}

func (h Hooks) AgentAction(ctx context.Context, action llm.AgentAction) {
	h.AgentActionFunc(ctx, action)
}

func (h Hooks) AgentFinish(ctx context.Context, finish llm.AgentFinish) {
	h.AgentFinishFunc(ctx, finish)
}

func (h Hooks) ChainStart(ctx context.Context, in map[string]any) {
	h.ChainStartFunc(ctx, in)
}

func (h Hooks) ChainEnd(ctx context.Context, out map[string]any) {
	h.ChainEndFunc(ctx, out)
}

func (h Hooks) ChainError(ctx context.Context, err error) {
	h.ChainErrorFunc(ctx, err)
}

func (h Hooks) ProviderStart(ctx context.Context, prompts []string) {
	h.ProviderStartFunc(ctx, prompts)
}

func (h Hooks) ProviderGenerateContentStart(ctx context.Context, messages []llm.Message) {
	h.ProviderGenerateContentStartFunc(ctx, messages)
}

func (h Hooks) ProviderGenerateContentEnd(ctx context.Context, res *llm.ContentResponse) {
	h.ProviderGenerateContentEndFunc(ctx, res)
}

func (h Hooks) ProviderError(ctx context.Context, err error) {
	h.ProviderErrorFunc(ctx, err)
}

func (h Hooks) RetrieverStart(ctx context.Context, query string) {
	h.RetrieverStartFunc(ctx, query)
}

func (h Hooks) RetrieverEnd(ctx context.Context, query string, docs []llm.Document) {
	h.RetrieverEndFunc(ctx, query, docs)
}

func (h Hooks) ToolStart(ctx context.Context, in string) {
	h.ToolStartFunc(ctx, in)
}

func (h Hooks) ToolEnd(ctx context.Context, out string) {
	h.ToolEndFunc(ctx, out)
}

func (h Hooks) ToolError(ctx context.Context, err error) {
	h.ToolErrorFunc(ctx, err)
}

func (h Hooks) AgentHooks() llm.AgentHooks {
	if h.AgentHooksFunc != nil {
		return h.AgentHooksFunc()
	}

	return h
}

func (h Hooks) ChainHooks() llm.ChainHooks {
	if h.ChainHooksFunc != nil {
		return h.ChainHooksFunc()
	}

	return h
}
