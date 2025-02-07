package mock

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Agent = Agent{}

type Agent struct {
	ToolsFunc      func() []llm.AgentTool
	InputKeysFunc  func() []string
	OutputKeysFunc func() []string
	PlanFunc       func(context.Context, []llm.AgentStep, map[string]string) ([]llm.AgentAction, *llm.AgentFinish, error)
}

func (a Agent) Tools() []llm.AgentTool {
	return a.ToolsFunc()
}

func (a Agent) InputKeys() []string {
	return a.InputKeysFunc()
}

func (a Agent) OutputKeys() []string {
	return a.OutputKeysFunc()
}

func (a Agent) Plan(ctx context.Context, steps []llm.AgentStep, in map[string]string) ([]llm.AgentAction, *llm.AgentFinish, error) {
	return a.PlanFunc(ctx, steps, in)
}
