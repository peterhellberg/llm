package agents

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/chains"
	"github.com/peterhellberg/llm/prompts"
)

const (
	conversationalFinalAnswerAction = "AI:"
)

// ConversationalAgent is a struct that represents an agent responsible for deciding
// what to do or give the final output if the task is finished given a set of inputs
// and previous steps taken.
//
// Other agents are often optimized for using tools to figure out the best response,
// which is not ideal in a conversational setting where you may want the agent to be
// able to chat with the user as well.
type ConversationalAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain llm.Chain
	// Tools is a list of the tools the agent can use.
	Tools []llm.AgentTool
	// Output key is the key where the final output is placed.
	OutputKey string
	// Hooks is the handler for callbacks.
	Hooks llm.Hooks
}

var _ llm.Agent = (*ConversationalAgent)(nil)

func NewConversationalAgent(provider llm.Provider, tools []llm.AgentTool, opts ...Option) *ConversationalAgent {
	options := conversationalDefaultOptions()

	for _, opt := range opts {
		opt(&options)
	}

	return &ConversationalAgent{
		Chain: chains.New(
			provider,
			options.getConversationalPrompt(tools),
			llm.ChainWithHooks(options.hooks),
		),
		Tools:     tools,
		OutputKey: options.outputKey,
		Hooks:     options.hooks,
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *ConversationalAgent) Plan(
	ctx context.Context,
	intermediateSteps []llm.AgentStep,
	inputs map[string]string,
) ([]llm.AgentAction, *llm.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))

	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

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

func (a *ConversationalAgent) InputKeys() []string {
	chainInputs := a.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *ConversationalAgent) OutputKeys() []string {
	return []string{a.OutputKey}
}

func (a *ConversationalAgent) AgentTools() []llm.AgentTool {
	return a.Tools
}

func constructScratchPad(steps []llm.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}

func (a *ConversationalAgent) parseOutput(output string) ([]llm.AgentAction, *llm.AgentFinish, error) {
	if strings.Contains(output, conversationalFinalAnswerAction) {
		splits := strings.Split(output, conversationalFinalAnswerAction)

		finishAction := &llm.AgentFinish{
			ReturnValues: map[string]any{
				a.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}

		return nil, finishAction, nil
	}

	r := regexp.MustCompile(`Action: (.*?)[\n]*Action Input: (.*)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", llm.ErrUnableToParseOutput, output)
	}

	return []llm.AgentAction{
		{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
	}, nil, nil
}

//go:embed prompts/conversational_prefix.txt
var defaultConversationalPrefix string

//go:embed prompts/conversational_format_instructions.txt
var defaultConversationalFormatInstructions string

//go:embed prompts/conversational_suffix.txt
var defaultConversationalSuffix string

func createConversationalPrompt(tools []llm.AgentTool, prefix, instructions, suffix string) prompts.Template {
	template := strings.Join([]string{prefix, instructions, suffix}, "\n\n")

	return prompts.Template{
		Template:       template,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad"},
		PartialVariables: map[string]any{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
			"history":           "",
		},
	}
}
