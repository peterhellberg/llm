package conversationchain

import "github.com/peterhellberg/llm"

const tmpl = `The following is a friendly conversation between a human and an AI. The AI is talkative and provides lots of specific details from its context. If the AI does not know the answer to a question, it truthfully says it does not know.

Current conversation:
{{.history}}
Human: {{.input}}
AI:`

func New(provider llm.Provider, opts ...llm.ChainOption) llm.Chain {
	opt := &llm.ChainOptions{}

	for _, o := range opts {
		o(opt)
	}

	return llm.NewChain(provider,
		llm.GoTemplate(tmpl, []string{
			"history",
			"input",
		}),
		llm.ChainWithMemory(opt.Memory),
		llm.ChainWithParser(opt.Parser),
		llm.ChainWithHooks(opt.Hooks),
	)
}
