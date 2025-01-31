package multiple

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Hooks = Hooks{}

// Hooks that combines multiple Hooks.
type Hooks struct {
	List []llm.Hooks
}

func (h Hooks) ProviderStart(ctx context.Context, prompts []string) {
	for _, l := range h.List {
		l.ProviderStart(ctx, prompts)
	}
}

func (h Hooks) ProviderGenerateContentStart(ctx context.Context, ms []llm.Message) {
	for _, l := range h.List {
		l.ProviderGenerateContentStart(ctx, ms)
	}
}

func (h Hooks) ProviderGenerateContentEnd(ctx context.Context, res *llm.ContentResponse) {
	for _, l := range h.List {
		l.ProviderGenerateContentEnd(ctx, res)
	}
}

func (h Hooks) ChainStart(ctx context.Context, inputs map[string]any) {
	for _, l := range h.List {
		l.ChainStart(ctx, inputs)
	}
}

func (h Hooks) ChainEnd(ctx context.Context, outputs map[string]any) {
	for _, l := range h.List {
		l.ChainEnd(ctx, outputs)
	}
}

func (h Hooks) ToolStart(ctx context.Context, input string) {
	for _, l := range h.List {
		l.ToolStart(ctx, input)
	}
}

func (h Hooks) ToolEnd(ctx context.Context, output string) {
	for _, l := range h.List {
		l.ToolEnd(ctx, output)
	}
}

func (h Hooks) AgentAction(ctx context.Context, action llm.AgentAction) {
	for _, l := range h.List {
		l.AgentAction(ctx, action)
	}
}

func (h Hooks) AgentFinish(ctx context.Context, finish llm.AgentFinish) {
	for _, l := range h.List {
		l.AgentFinish(ctx, finish)
	}
}

func (h Hooks) RetrieverStart(ctx context.Context, query string) {
	for _, l := range h.List {
		l.RetrieverStart(ctx, query)
	}
}

func (h Hooks) RetrieverEnd(ctx context.Context, query string, documents []llm.Document) {
	for _, l := range h.List {
		l.RetrieverEnd(ctx, query, documents)
	}
}

func (h Hooks) StreamingFunc(ctx context.Context, chunk []byte) {
	for _, l := range h.List {
		l.StreamingFunc(ctx, chunk)
	}
}

func (h Hooks) ChainError(ctx context.Context, err error) {
	for _, l := range h.List {
		l.ChainError(ctx, err)
	}
}

func (h Hooks) ProviderError(ctx context.Context, err error) {
	for _, l := range h.List {
		l.ProviderError(ctx, err)
	}
}

func (h Hooks) ToolError(ctx context.Context, err error) {
	for _, l := range h.List {
		l.ToolError(ctx, err)
	}
}
