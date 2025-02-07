package messageformatters

import "github.com/peterhellberg/llm"

var _ llm.MessageFormatter = AI{}

// NewAI creates a new AI message prompt formatter.
func NewAI(content string, variables []string) AI {
	return AI{
		Template: llm.GoTemplate(content, variables),
	}
}

// AI is a message formatter that returns an AI message.
type AI struct {
	Template llm.Template
}

// FormatMessages formats the message with the values given.
func (a AI) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	text, err := a.Template.FormatString(values)

	return []llm.ChatMessage{
		llm.AIChatMessage{Content: text},
	}, err
}

// InputVariables returns the input variables the prompt expects.
func (a AI) InputVariables() []string {
	return a.Template.InputVariables()
}
