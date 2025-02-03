package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Memory = Memory{}

type Memory struct {
	ClearFunc         func(context.Context) error
	LoadVariablesFunc func(context.Context, map[string]any) (map[string]any, error)
	MemoryKeyFunc     func(context.Context) string
	SaveContextFunc   func(context.Context, map[string]any, map[string]any) error
	VariablesFunc     func(context.Context) []string
}

func (m Memory) Clear(ctx context.Context) error {
	return m.ClearFunc(ctx)
}

func (m Memory) LoadVariables(ctx context.Context, in map[string]any) (map[string]any, error) {
	return m.LoadVariablesFunc(ctx, in)
}

func (m Memory) MemoryKey(ctx context.Context) string {
	return m.MemoryKeyFunc(ctx)
}

func (m Memory) SaveContext(ctx context.Context, in map[string]any, out map[string]any) error {
	return m.SaveContextFunc(ctx, in, out)
}

func (m Memory) Variables(ctx context.Context) []string {
	return m.VariablesFunc(ctx)
}
