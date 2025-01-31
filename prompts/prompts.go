package prompts

import "github.com/peterhellberg/llm"

var (
	_ llm.StringFormatter = Template{}
	_ llm.Prompter        = Template{}
)

// Template contains common fields for all prompt templates.
type Template struct {
	// Template is the prompt template.
	Template string

	// A list of variable names the prompt template expects.
	InputVariables []string

	// TemplateFormat is the format of the prompt template.
	TemplateFormat TemplateFormat

	// OutputParser is a function that parses the output of the prompt template.
	OutputParser llm.Parser[any]

	// PartialVariables represents a map of variable names to values or functions
	// that return values. If the value is a function, it will be called when the
	// prompt template is rendered.
	PartialVariables map[string]any
}

// NewTemplate returns a new prompt template.
func NewTemplate(template string, inputVars []string) Template {
	return Template{
		Template:       template,
		InputVariables: inputVars,
		TemplateFormat: TemplateFormatGoTemplate,
	}
}

// Format formats the prompt template and returns a string value.
func (p Template) FormatString(values map[string]any) (string, error) {
	resolvedValues, err := llm.ResolvePartialValues(p.PartialVariables, values)
	if err != nil {
		return "", err
	}

	return RenderTemplate(p.Template, p.TemplateFormat, resolvedValues)
}

// FormatPrompt formats the prompt template and returns a string prompt value.
func (p Template) FormatPrompt(values map[string]any) (llm.Prompt, error) {
	f, err := p.FormatString(values)
	if err != nil {
		return nil, err
	}

	return String(f), nil
}

// GetInputVariables returns the input variables the prompt expect.
func (p Template) GetInputVariables() []string {
	return p.InputVariables
}
