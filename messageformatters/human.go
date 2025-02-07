package messageformatters

import "github.com/peterhellberg/llm"

var _ llm.MessageFormatter = Human{}

// NewHuman creates a new human message prompt formatter.
func NewHuman(content string, variables []string) Human {
	return Human{
		Template: llm.GoTemplate(content, variables),
	}
}

// Human is a message formatter that returns a human message.
type Human struct {
	Template llm.Template
}

// FormatMessages formats the message with the values given.
func (h Human) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	text, err := h.Template.FormatString(values)

	return []llm.ChatMessage{
		llm.HumanChatMessage{Content: text},
	}, err
}

// InputVariables returns the input variables the prompt expects.
func (h Human) InputVariables() []string {
	return h.Template.InputVariables()
}
