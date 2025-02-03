package llm

import (
	"html/template"
	"strings"
)

// GoTemplate constructs a Template that renders using text/template
func GoTemplate(content string, variables []string) Template {
	return NewTemplate(content, variables, RendererFunc(
		func(tmpl string, values map[string]any) (string, error) {
			parsed, err := template.New("template").Option("missingkey=error").Parse(tmpl)
			if err != nil {
				return "", err
			}

			sb := new(strings.Builder)

			if err := parsed.Execute(sb, values); err != nil {
				return "", err
			}

			return sb.String(), nil
		},
	))
}
