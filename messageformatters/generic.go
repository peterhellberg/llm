package messageformatters

import "github.com/peterhellberg/llm"

var _ llm.MessageFormatter = Generic{}

// NewGeneric creates a new generic message prompt formatter.
func NewGeneric(role, content string, variables []string) Generic {
	return Generic{
		Template: llm.GoTemplate(content, variables),
		Role:     role,
	}
}

// Generic is a message formatter that returns message with the specified speaker.
type Generic struct {
	Template llm.Template
	Role     string
}

// FormatMessages formats the message with the values given.
func (g Generic) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	text, err := g.Template.FormatString(values)

	return []llm.ChatMessage{
		llm.GenericChatMessage{Content: text, Role: g.Role},
	}, err
}

// InputVariables returns the input variables the prompt expects.
func (g Generic) InputVariables() []string {
	return g.Template.InputVariables()
}
