package prompts

import "github.com/peterhellberg/llm"

var (
	_ llm.StringFormatter = Template{}
	_ llm.PromptFormatter = Template{}
)

type Formatter interface {
	RenderTemplate(tmpl string, values map[string]any) (string, error)
}

// Template contains common fields for all prompt templates.
type Template struct {
	// Content of the prompt template.
	Content string

	// Variables is a list of variable names the prompt template expects.
	Variables []string

	// Parser is a function that parses the output of the prompt template.
	Parser llm.Parser[any]

	// PartialVariables represents a map of variable names to values or functions
	// that return values. If the value is a function, it will be called when the
	// prompt template is rendered.
	PartialVariables map[string]any

	// Formatter used to generate strings based on the prompt template.
	Formatter Formatter
}

// NewTemplate returns a new prompt template.
func NewTemplate(content string, variables []string, formatter Formatter) Template {
	return Template{
		Content:   content,
		Variables: variables,
		Formatter: formatter,
	}
}

// Format formats the prompt template and returns a string value.
func (p Template) FormatString(values map[string]any) (string, error) {
	resolvedValues, err := llm.ResolvePartialValues(p.PartialVariables, values)
	if err != nil {
		return "", err
	}

	return p.Formatter.RenderTemplate(p.Content, resolvedValues)
}

// FormatPrompt formats the prompt template and returns a string prompt value.
func (p Template) FormatPrompt(values map[string]any) (llm.Prompt, error) {
	f, err := p.FormatString(values)
	if err != nil {
		return nil, err
	}

	return String(f), nil
}

// InputVariables returns the input variables the prompt expect.
func (p Template) InputVariables() []string {
	return p.Variables
}
