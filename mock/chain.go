package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var (
	_ llm.Chain       = Chain{}
	_ llm.ChainHooker = Chain{}
)

type Chain struct {
	CallFunc       func(context.Context, map[string]any, ...llm.ChainOption) (map[string]any, error)
	MemoryFunc     func() llm.Memory
	InputKeysFunc  func() []string
	OutputKeysFunc func() []string
	ChainHooksFunc func() llm.ChainHooks
}

func (c Chain) Call(ctx context.Context, values map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	return c.CallFunc(ctx, values, options...)
}

func (c Chain) Memory() llm.Memory {
	return c.MemoryFunc()
}

func (c Chain) InputKeys() []string {
	return c.InputKeysFunc()
}

func (c Chain) OutputKeys() []string {
	return c.OutputKeysFunc()
}

func (c Chain) ChainHooks() llm.ChainHooks {
	return c.ChainHooksFunc()
}
