package gotemplate

import (
	"strings"
	"text/template"
)

type Formatter struct{}

func (Formatter) RenderTemplate(tmpl string, values map[string]any) (string, error) {
	parsedTmpl, err := template.New("template").
		Option("missingkey=error").
		Parse(tmpl)
	if err != nil {
		return "", err
	}

	sb := new(strings.Builder)

	if err := parsedTmpl.Execute(sb, values); err != nil {
		return "", err
	}

	return sb.String(), nil
}
