package llm

import (
	"context"
)

// Agent is the interface all agents must implement.
type Agent interface {
	// Plan Given an input and previous steps decide what to do next. Returns either actions or a finish.
	Plan(ctx context.Context, steps []AgentStep, inputs map[string]string) ([]AgentAction, *AgentFinish, error)
	InputKeys() []string
	OutputKeys() []string
	AgentTools() []AgentTool
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

// ParserErrorHandler is the struct used to handle parse errors from the agent in the executor. If
// an executor have a ParserErrorHandler, parsing errors will be formatted using the formatter
// function and added as an observation. In the next executor step the agent will then have the
// possibility to fix the error.
type ParserErrorHandler struct {
	// The formatter function can be used to format the parsing error. If nil the error will be given
	// as an observation directly.
	Formatter func(err string) string
}

// NewParserErrorHandler creates a new parser error handler.
func NewParserErrorHandler(formatFunc func(string) string) *ParserErrorHandler {
	return &ParserErrorHandler{
		Formatter: formatFunc,
	}
}
