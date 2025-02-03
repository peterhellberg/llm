package chains

import (
	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/parsers"
	"github.com/peterhellberg/llm/prompts"
	"github.com/peterhellberg/llm/prompts/formatters/gotemplate"
)

const conversationTemplate = `The following is a friendly conversation between a human and an AI. The AI is talkative and provides lots of specific details from its context. If the AI does not know the answer to a question, it truthfully says it does not know.

Current conversation:
{{.history}}
Human: {{.input}}
AI:`

func Conversation(provider llm.Provider, memory llm.Memory, opts ...llm.ChainOption) *Chain {
	opt := &llm.ChainOptions{}

	for _, o := range opts {
		o(opt)
	}

	return &Chain{
		Hooks: opt.Hooks,

		prompter: prompts.NewTemplate(
			conversationTemplate,
			[]string{"history", "input"},
			gotemplate.Formatter{},
		),
		provider: provider,
		memory:   memory,
		parser:   parsers.Empty{},

		outputKey: defaultOutputKey,
	}
}
