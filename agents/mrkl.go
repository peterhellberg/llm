package agents

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/chains/llmchain"
)

const (
	finalAnswerAction = "Final Answer:"
	defaultOutputKey  = "output"
)

// OneShotZeroAgent is a struct that represents an agent responsible for deciding
// what to do or give the final output if the task is finished given a set of inputs
// and previous steps taken.
//
// This agent is optimized to be used with LLMs.
type OneShotZeroAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain llm.Chain
	// Tools is a list of the tools the agent can use.
	Tools []llm.AgentTool
	// Output key is the key where the final output is placed.
	OutputKey string
	// Hooks for the Agent.
	Hooks llm.AgentHooks
}

var _ llm.Agent = (*OneShotZeroAgent)(nil)

// NewOneShotAgent creates a new OneShotZeroAgent with the given LLM model, tools,
// and options. It returns a pointer to the created agent. The opts parameter
// represents the options for the agent.
func NewOneShotAgent(provider llm.Provider, tools []llm.AgentTool, opts ...Option) *OneShotZeroAgent {
	options := mrklDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &OneShotZeroAgent{
		Chain: llmchain.New(
			provider,
			options.getMrklPrompt(tools),
			llm.ChainWithHooks(options.hooks),
		),
		Tools:     tools,
		OutputKey: options.outputKey,
		Hooks:     options.hooks,
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *OneShotZeroAgent) Plan(
	ctx context.Context,
	intermediateSteps []llm.AgentStep,
	inputs map[string]string,
) ([]llm.AgentAction, *llm.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))

	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructMrklScratchPad(intermediateSteps)
	fullInputs["today"] = time.Now().Format("January 02, 2006")

	var stream func(ctx context.Context, chunk []byte) error

	if a.Hooks != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			a.Hooks.StreamingFunc(ctx, chunk)
			return nil
		}
	}

	output, err := llm.ChainPredict(
		ctx,
		a.Chain,
		fullInputs,
		llm.ChainWithStopWords([]string{"\nObservation:", "\n\tObservation:"}),
		llm.ChainWithStreamingFunc(stream),
	)
	if err != nil {
		return nil, nil, err
	}

	return a.parseOutput(output)
}

func (a *OneShotZeroAgent) InputKeys() []string {
	chainInputs := a.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" || v == "today" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *OneShotZeroAgent) OutputKeys() []string {
	return []string{a.OutputKey}
}

func (a *OneShotZeroAgent) AgentTools() []llm.AgentTool {
	return a.Tools
}

func constructMrklScratchPad(steps []llm.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += "\n" + step.Action.Log
			scratchPad += "\nObservation: " + step.Observation + "\n"
		}
	}

	return scratchPad
}

func (a *OneShotZeroAgent) parseOutput(output string) ([]llm.AgentAction, *llm.AgentFinish, error) {
	if strings.Contains(output, finalAnswerAction) {
		splits := strings.Split(output, finalAnswerAction)

		return nil, &llm.AgentFinish{
			ReturnValues: map[string]any{
				a.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}, nil
	}

	r := regexp.MustCompile(`Action:\s*(.+)\s*Action Input:\s(?s)*(.+)`)

	matches := r.FindStringSubmatch(output)

	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", llm.ErrUnableToParseOutput, output)
	}

	return []llm.AgentAction{
		{
			Tool:      strings.TrimSpace(matches[1]),
			ToolInput: strings.TrimSpace(matches[2]),
			Log:       output,
		},
	}, nil, nil
}
