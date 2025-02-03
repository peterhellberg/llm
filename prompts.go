package llm

// StringFormatter is an interface for formatting a map of values into a string.
type StringFormatter interface {
	FormatString(values map[string]any) (string, error)
}

// MessageFormatter is an interface for formatting a map of values into a list
// of messages.
type MessageFormatter interface {
	FormatMessages(values map[string]any) ([]ChatMessage, error)
	InputVariables() []string
}

// PromptFormatter is an interface for formatting a map of values into a prompt.
type PromptFormatter interface {
	FormatPrompt(values map[string]any) (Prompt, error)
	InputVariables() []string
}

// Prompt is the interface that all prompt values must implement.
type Prompt interface {
	String() string
	Messages() []ChatMessage
}
