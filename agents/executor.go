package agents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/peterhellberg/llm"
)

var (
	_ llm.Chain       = &Executor{}
	_ llm.AgentHooker = &Executor{}
)

const _intermediateStepsOutputKey = "intermediateSteps"

// Executor is the chain responsible for running agents.
type Executor struct {
	Agent        llm.Agent
	Memory       llm.Memory
	Hooks        llm.AgentHooks
	ErrorHandler *llm.ParserErrorHandler

	MaxIterations           int
	ReturnIntermediateSteps bool
}

// NewExecutor creates a new agent executor with an agent and the tools the agent can use.
func NewExecutor(agent llm.Agent, opts ...Option) *Executor {
	options := executorDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Executor{
		Agent:                   agent,
		Memory:                  options.memory,
		MaxIterations:           options.maxIterations,
		ReturnIntermediateSteps: options.returnIntermediateSteps,
		Hooks:                   options.hooks,
		ErrorHandler:            options.errorHandler,
	}
}

func (e *Executor) Call(ctx context.Context, inputValues map[string]any, _ ...llm.ChainOption) (map[string]any, error) {
	inputs, err := inputsToString(inputValues)
	if err != nil {
		return nil, err
	}
	nameToTool := getNameToTool(e.Agent.AgentTools())

	steps := make([]llm.AgentStep, 0)
	for i := 0; i < e.MaxIterations; i++ {
		var finish map[string]any
		steps, finish, err = e.doIteration(ctx, steps, nameToTool, inputs)
		if finish != nil || err != nil {
			return finish, err
		}
	}

	if e.Hooks != nil {
		e.Hooks.AgentFinish(ctx, llm.AgentFinish{
			ReturnValues: map[string]any{"output": llm.ErrNotFinished.Error()},
		})
	}
	return e.getReturn(
		&llm.AgentFinish{ReturnValues: make(map[string]any)},
		steps,
	), llm.ErrNotFinished
}

func (e *Executor) doIteration(
	ctx context.Context,
	steps []llm.AgentStep,
	nameToTool map[string]llm.AgentTool,
	inputs map[string]string,
) ([]llm.AgentStep, map[string]any, error) {
	actions, finish, err := e.Agent.Plan(ctx, steps, inputs)

	if errors.Is(err, llm.ErrUnableToParseOutput) && e.ErrorHandler != nil {
		formattedObservation := err.Error()
		if e.ErrorHandler.Formatter != nil {
			formattedObservation = e.ErrorHandler.Formatter(formattedObservation)
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
		if e.Hooks != nil {
			e.Hooks.AgentFinish(ctx, *finish)
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

func (e *Executor) doAction(
	ctx context.Context,
	steps []llm.AgentStep,
	nameToTool map[string]llm.AgentTool,
	action llm.AgentAction,
) ([]llm.AgentStep, error) {
	if e.Hooks != nil {
		e.Hooks.AgentAction(ctx, action)
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
	if e.ReturnIntermediateSteps {
		finish.ReturnValues[_intermediateStepsOutputKey] = steps
	}

	return finish.ReturnValues
}

// GetInputKeys gets the input keys the agent of the executor expects.
// Often "input".
func (e *Executor) GetInputKeys() []string {
	return e.Agent.InputKeys()
}

// GetOutputKeys gets the output keys the agent of the executor returns.
func (e *Executor) GetOutputKeys() []string {
	return e.Agent.OutputKeys()
}

func (e *Executor) GetMemory() llm.Memory {
	return e.Memory
}

func (e *Executor) AgentHooks() llm.AgentHooks {
	return e.Hooks
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
