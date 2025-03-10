package agents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/peterhellberg/llm"
)

const (
	defaultMaxIterations       = 5
	defaultOutputKey           = "output"
	intermediateStepsOutputKey = "intermediateSteps"
)

func executorDefaultOptions() Options {
	return Options{
		maxIterations: defaultMaxIterations,
		outputKey:     defaultOutputKey,
		memory:        llm.EmptyMemory{},
	}
}

// Executor is the chain responsible for running agents.
type Executor struct {
	agent        llm.Agent
	memory       llm.Memory
	hooks        llm.AgentHooks
	errorHandler *llm.ParserErrorHandler

	maxIterations           int
	returnIntermediateSteps bool
}

var (
	_ llm.Chain       = &Executor{}
	_ llm.AgentHooker = &Executor{}
)

// NewExecutor creates a new agent executor with an agent and the tools the agent can use.
func NewExecutor(agent llm.Agent, opts ...Option) *Executor {
	options := executorDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Executor{
		agent:                   agent,
		memory:                  options.memory,
		maxIterations:           options.maxIterations,
		returnIntermediateSteps: options.returnIntermediateSteps,
		hooks:                   options.hooks,
		errorHandler:            options.errorHandler,
	}
}

func (e *Executor) Call(ctx context.Context, inputValues map[string]any, _ ...llm.ChainOption) (map[string]any, error) {
	inputs, err := inputsToString(inputValues)
	if err != nil {
		return nil, err
	}

	nameToTool := getNameToTool(e.agent.Tools())

	steps := make([]llm.AgentStep, 0)

	for i := 0; i < e.maxIterations; i++ {
		var finish map[string]any

		steps, finish, err = e.doIteration(ctx, steps, nameToTool, inputs)
		if finish != nil || err != nil {
			return finish, err
		}
	}

	if e.hooks != nil {
		e.hooks.AgentFinish(ctx, llm.AgentFinish{
			ReturnValues: map[string]any{"output": llm.ErrNotFinished.Error()},
		})
	}

	return e.getReturn(
		&llm.AgentFinish{ReturnValues: make(map[string]any)}, steps,
	), llm.ErrNotFinished
}

func (e *Executor) doIteration(ctx context.Context, steps []llm.AgentStep, nameToTool map[string]llm.AgentTool, inputs map[string]string) ([]llm.AgentStep, map[string]any, error) {
	actions, finish, err := e.agent.Plan(ctx, steps, inputs)
	if errors.Is(err, llm.ErrUnableToParseOutput) && e.errorHandler != nil {
		formattedObservation := err.Error()

		if e.errorHandler.Formatter != nil {
			formattedObservation = e.errorHandler.Formatter(formattedObservation)
		}

		steps = append(steps, llm.AgentStep{
			Observation: formattedObservation,
		})

		return steps, nil, nil
	}

	if err != nil {
		return steps, nil, err
	}

	if len(actions) == 0 && finish == nil {
		return steps, nil, llm.ErrAgentNoReturn
	}

	if finish != nil {
		if e.hooks != nil {
			e.hooks.AgentFinish(ctx, *finish)
		}

		return steps, e.getReturn(finish, steps), nil
	}

	for _, action := range actions {
		steps, err = e.doAction(ctx, steps, nameToTool, action)
		if err != nil {
			return steps, nil, err
		}
	}

	return steps, nil, nil
}

func (e *Executor) doAction(ctx context.Context, steps []llm.AgentStep, nameToTool map[string]llm.AgentTool, action llm.AgentAction) ([]llm.AgentStep, error) {
	if e.hooks != nil {
		e.hooks.AgentAction(ctx, action)
	}

	tool, ok := nameToTool[strings.ToUpper(action.Tool)]
	if !ok {
		return append(steps, llm.AgentStep{
			Action:      action,
			Observation: fmt.Sprintf("%s is not a valid tool, try another one", action.Tool),
		}), nil
	}

	observation, err := tool.Call(ctx, action.ToolInput)
	if err != nil {
		return nil, err
	}

	return append(steps, llm.AgentStep{
		Action:      action,
		Observation: observation,
	}), nil
}

func (e *Executor) getReturn(finish *llm.AgentFinish, steps []llm.AgentStep) map[string]any {
	if e.returnIntermediateSteps {
		finish.ReturnValues[intermediateStepsOutputKey] = steps
	}

	return finish.ReturnValues
}

// GetInputKeys gets the input keys the agent of the executor expects.
// Often "input".
func (e *Executor) InputKeys() []string {
	return e.agent.InputKeys()
}

// OutputKeys gets the output keys the agent of the executor returns.
func (e *Executor) OutputKeys() []string {
	return e.agent.OutputKeys()
}

func (e *Executor) Memory() llm.Memory {
	return e.memory
}

func (e *Executor) AgentHooks() llm.AgentHooks {
	return e.hooks
}

func inputsToString(inputValues map[string]any) (map[string]string, error) {
	inputs := make(map[string]string, len(inputValues))

	for key, value := range inputValues {
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s", llm.ErrExecutorInputNotString, key)
		}

		inputs[key] = valueStr
	}

	return inputs, nil
}

func getNameToTool(t []llm.AgentTool) map[string]llm.AgentTool {
	if len(t) == 0 {
		return nil
	}

	nameToTool := make(map[string]llm.AgentTool, len(t))

	for _, tool := range t {
		nameToTool[strings.ToUpper(tool.Name())] = tool
	}

	return nameToTool
}
