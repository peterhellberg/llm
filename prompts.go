package llm

// StringFormatter is an interface for formatting a map of values into a string.
type StringFormatter interface {
	FormatString(values map[string]any) (string, error)
}

// MessageFormatter is an interface for formatting a map of values into a list
// of messages.
type MessageFormatter interface {
	FormatMessages(values map[string]any) ([]ChatMessage, error)
	GetInputVariables() []string
}

// Prompter is an interface for formatting a map of values into a prompt.
type Prompter interface {
	FormatPrompt(values map[string]any) (Prompt, error)
	GetInputVariables() []string
}

// Prompt is the interface that all prompt values must implement.
type Prompt interface {
	String() string
	Messages() []ChatMessage
}
