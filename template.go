package llm

var (
	_ StringFormatter = Template{}
	_ PromptFormatter = Template{}
)

func NewTemplate(content string, variables []string, renderer Renderer, options ...TemplateOption) Template {
	t := Template{
		Content:   content,
		Variables: variables,
		Renderer:  renderer,
	}

	for _, o := range options {
		o(&t)
	}

	return t
}

type TemplateOption func(*Template)

func TemplateWithParser(parser Parser[any]) TemplateOption {
	return func(t *Template) {
		t.Parser = parser
	}
}

func TemplateWithPartialVariables(partialVariables map[string]any) TemplateOption {
	return func(t *Template) {
		t.PartialVariables = partialVariables
	}
}

// Template contains common fields for all prompt templates.
type Template struct {
	// Content of the prompt template.
	Content string

	// Variables is a list of variable names the prompt template expects.
	Variables []string

	// Renderer used to generate strings based on the prompt template.
	Renderer

	// Parser is a function that parses the output of the prompt template.
	Parser Parser[any]

	// PartialVariables represents a map of variable names to values or functions
	// that return values. If the value is a function, it will be called when the
	// prompt template is rendered.
	PartialVariables map[string]any
}

// Format formats the prompt template and returns a string value.
func (p Template) FormatString(values map[string]any) (string, error) {
	resolvedValues, err := ResolvePartialValues(p.PartialVariables, values)
	if err != nil {
		return "", err
	}

	return p.RenderTemplate(p.Content, resolvedValues)
}

// FormatPrompt formats the prompt template and returns a string prompt value.
func (p Template) FormatPrompt(values map[string]any) (Prompt, error) {
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
