package llm

import "github.com/peterhellberg/llm/internal/fstring"

// FStringTemplate constructs a TEmplate that renders using internal/fstring
func FStringTemplate(content string, variables []string) Template {
	renderer := RendererFunc(
		func(tmpl string, values map[string]any) (string, error) {
			return fstring.Format(tmpl, values)
		},
	)

	return NewTemplate(content, variables, renderer)
}
