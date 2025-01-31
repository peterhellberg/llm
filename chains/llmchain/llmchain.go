package llmchain

import (
	"context"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/memory"
	"github.com/peterhellberg/llm/parsers"
)

var (
	_ llm.Chain       = &Chain{}
	_ llm.ChainHooker = &Chain{}
)

const defaultOutputKey = "text"

type Chain struct {
	Hooks llm.ChainHooks

	prompter llm.Prompter
	provider llm.Provider
	parser   llm.Parser[any]
	memory   llm.Memory

	OutputKey string
}

// New chain with an LLM and a prompt.
func New(provider llm.Provider, prompter llm.Prompter, opts ...llm.ChainOption) *Chain {
	opt := &llm.ChainOptions{}

	for _, o := range opts {
		o(opt)
	}

	return &Chain{
		Hooks: opt.Hooks,

		prompter: prompter,
		provider: provider,
		parser:   parsers.Empty{},
		memory:   memory.Empty{},

		OutputKey: defaultOutputKey,
	}
}

// Call formats the prompts with the input values, generates using the llm, and parses
// the output from the llm with the output parser. This function should not be called
// directly, use rather the Call or Run function if the prompt only requires one input value.
func (c Chain) Call(ctx context.Context, values map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	prompt, err := c.prompter.FormatPrompt(values)
	if err != nil {
		return nil, err
	}

	result, err := llm.Call(ctx, c.provider, prompt.String(), llm.ChainToContentOptions(options...)...)
	if err != nil {
		return nil, err
	}

	finalOutput, err := c.parser.ParseWithPrompt(result, prompt)
	if err != nil {
		return nil, err
	}

	return map[string]any{c.OutputKey: finalOutput}, nil
}

// GetMemory returns the memory.
func (c Chain) GetMemory() llm.Memory {
	return c.memory
}

func (c Chain) ChainHooks() llm.ChainHooks {
	return c.Hooks
}

// GetInputKeys returns the expected input keys.
func (c Chain) GetInputKeys() []string {
	return append([]string{}, c.prompter.GetInputVariables()...)
}

// GetOutputKeys returns the output keys the chain will return.
func (c Chain) GetOutputKeys() []string {
	return []string{c.OutputKey}
}
