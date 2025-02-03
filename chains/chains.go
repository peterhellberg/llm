package chains

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

	prompter llm.PromptFormatter
	provider llm.Provider
	parser   llm.Parser[any]
	memory   llm.Memory

	outputKey string
}

// New chain with a LLM provider and a prompt.
func New(provider llm.Provider, prompter llm.PromptFormatter, opts ...llm.ChainOption) *Chain {
	opt := &llm.ChainOptions{}

	for _, o := range opts {
		o(opt)
	}

	{
		if opt.Parser == nil {
			opt.Parser = parsers.Empty{}
		}

		if opt.Memory == nil {
			opt.Memory = memory.Empty{}
		}

		if opt.OutputKey == "" {
			opt.OutputKey = defaultOutputKey
		}
	}

	return &Chain{
		provider: provider,
		prompter: prompter,

		parser:    opt.Parser,
		memory:    opt.Memory,
		outputKey: opt.OutputKey,

		Hooks: opt.Hooks,
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

	return map[string]any{c.outputKey: finalOutput}, nil
}

// Memory returns the memory.
func (c Chain) Memory() llm.Memory {
	return c.memory
}

// ChainHooks returns the hooks for the chain.
func (c Chain) ChainHooks() llm.ChainHooks {
	return c.Hooks
}

// InputKeys returns the expected input keys.
func (c Chain) InputKeys() []string {
	return append([]string{}, c.prompter.InputVariables()...)
}

// OutputKeys returns the output keys the chain will return.
func (c Chain) OutputKeys() []string {
	return []string{c.outputKey}
}
