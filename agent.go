package llm

import "context"

// Agent is the interface all agents must implement.
type Agent interface {
	// Plan Given an input and previous steps decide what to do next. Returns either actions or a finish.
	Plan(ctx context.Context, steps []AgentStep, inputs map[string]string) ([]AgentAction, *AgentFinish, error)
	InputKeys() []string
	OutputKeys() []string
	Tools() []AgentTool
}

// AgentAction is the agent's action to take.
type AgentAction struct {
	Tool      string
	ToolInput string
	Log       string
	ToolID    string
}

// AgentStep is a step of the agent.
type AgentStep struct {
	Action      AgentAction
	Observation string
}

// AgentFinish is the agent's return value.
type AgentFinish struct {
	ReturnValues map[string]any
	Log          string
}

// AgentTool is a tool for the LLM agent to interact with different applications.
type AgentTool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string) (string, error)
}
