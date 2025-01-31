package prompts

import "github.com/peterhellberg/llm"

var _ llm.Prompt = String("")

// String is a prompt value that is a string.
type String string

func (s String) String() string {
	return string(s)
}

// Messages returns a single-element ChatMessage slice.
func (s String) Messages() []llm.ChatMessage {
	return []llm.ChatMessage{
		llm.HumanChatMessage{
			Content: string(s),
		},
	}
}
