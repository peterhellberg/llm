package messageformatters

import (
	"fmt"

	"github.com/peterhellberg/llm"
)

var _ llm.MessageFormatter = Placeholder{}

// NewPlaceholder creates a new placeholder message prompt formatter.
func NewPlaceholder(variable string) Placeholder {
	return Placeholder{
		Variable: variable,
	}
}

// Placeholder is a message formatter that returns a placeholder message.
type Placeholder struct {
	Variable string
}

// FormatMessages formats the messages from the values by variable name.
func (p Placeholder) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	value, ok := values[p.Variable]
	if !ok {
		return nil, fmt.Errorf(
			"%w: %s should be a list of chat messages",
			llm.ErrNeedChatMessageList, p.Variable)
	}

	base, ok := value.([]llm.ChatMessage)
	if !ok {
		return nil, fmt.Errorf(
			"%w: %s should be a list of chat messages",
			llm.ErrNeedChatMessageList, p.Variable)
	}

	return base, nil
}

// InputVariables returns the input variables the prompt expect.
func (p Placeholder) InputVariables() []string {
	return []string{p.Variable}
}
