package chatprompt

import "github.com/peterhellberg/llm"

var (
	_ llm.StringFormatter  = Template{}
	_ llm.MessageFormatter = Template{}
	_ llm.PromptFormatter  = Template{}
)

// Template is a prompt template for chat messages.
type Template struct {
	// Messages is the list of the messages to be formatted.
	Messages []llm.MessageFormatter

	// PartialVariables represents a map of variable names to values or functions
	// that return values. If the value is a function, it will be called when the
	// prompt template is rendered.
	PartialVariables map[string]any
}

// FormatPrompt formats the messages into a chat prompt value.
func (p Template) FormatPrompt(values map[string]any) (llm.Prompt, error) {
	resolvedValues, err := llm.ResolvePartialValues(p.PartialVariables, values)
	if err != nil {
		return nil, err
	}

	formattedMessages := make([]llm.ChatMessage, 0, len(p.Messages))

	for _, m := range p.Messages {
		curFormattedMessages, err := m.FormatMessages(resolvedValues)
		if err != nil {
			return nil, err
		}

		formattedMessages = append(formattedMessages, curFormattedMessages...)
	}

	return Value(formattedMessages), nil
}

// FormatString formats the messages with values given and returns the messages as a string.
func (p Template) FormatString(values map[string]any) (string, error) {
	promptValue, err := p.FormatPrompt(values)

	return promptValue.String(), err
}

// FormatMessages formats the messages with the values and returns the formatted messages.
func (p Template) FormatMessages(values map[string]any) ([]llm.ChatMessage, error) {
	promptValue, err := p.FormatPrompt(values)

	if promptValue == nil {
		return nil, err
	}

	return promptValue.Messages(), err
}

// InputVariables returns the input variables the prompt expect.
func (p Template) InputVariables() []string {
	inputVariablesMap := make(map[string]bool, 0)

	for _, msg := range p.Messages {
		for _, variable := range msg.InputVariables() {
			inputVariablesMap[variable] = true
		}
	}

	inputVariables := make([]string, 0, len(inputVariablesMap))

	for variable := range inputVariablesMap {
		inputVariables = append(inputVariables, variable)
	}

	return inputVariables
}

// NewTemplate creates a new chat prompt template from a list of message formatters.
func NewTemplate(messages []llm.MessageFormatter) Template {
	return Template{
		Messages: messages,
	}
}
