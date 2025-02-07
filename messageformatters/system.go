package messageformatters

import "github.com/peterhellberg/llm"

var _ llm.MessageFormatter = System{}

// NewSystem creates a new system message prompt formatter.
func NewSystem(content string, inputVariables []string) System {
	return System{
		Template: llm.GoTemplate(content, inputVariables),
	}
}

// System is a message formatter that returns a system message.
type System struct {
	Template llm.Template
}

// FormatMessages formats the message with the values given.
func (s System) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	text, err := s.Template.FormatString(values)

	return []llm.ChatMessage{
		llm.SystemChatMessage{Content: text},
	}, err
}

// InputVariables returns the input variables the prompt expects.
func (s System) InputVariables() []string {
	return s.Template.InputVariables()
}
