package prompts

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"text/template"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/prompts/internal/fstring"
)

// TemplateFormat is the format of the template.
type TemplateFormat string

const (
	// TemplateFormatGoTemplate is the format for go-template.
	TemplateFormatGoTemplate TemplateFormat = "go-template"

	// TemplateFormatFString is the format for f-string.
	TemplateFormatFString TemplateFormat = "f-string"
)

// interpolator is the function that interpolates the given template with the given values.
type interpolator func(template string, values map[string]any) (string, error)

// defaultFormatterMapping is the default mapping of TemplateFormat to interpolator.
var defaultFormatterMapping = map[TemplateFormat]interpolator{
	TemplateFormatGoTemplate: interpolateGoTemplate,
	TemplateFormatFString:    fstring.Format,
}

// interpolateGoTemplate interpolates the given template with the given values by using
// text/template.
func interpolateGoTemplate(tmpl string, values map[string]any) (string, error) {
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

func newInvalidTemplateError(gotTemplateFormat TemplateFormat) error {
	keys := maps.Keys(defaultFormatterMapping)

	formats := slices.Collect(keys)

	slices.Sort(formats)

	return fmt.Errorf("%w, got: %s, should be one of %s",
		llm.ErrInvalidTemplateFormat,
		gotTemplateFormat,
		formats,
	)
}

// CheckValidTemplate checks if the template is valid through checking whether the given
// TemplateFormat is available and whether the template can be rendered.
func CheckValidTemplate(template string, templateFormat TemplateFormat, inputVariables []string) error {
	_, ok := defaultFormatterMapping[templateFormat]
	if !ok {
		return newInvalidTemplateError(templateFormat)
	}

	dummyInputs := make(map[string]any, len(inputVariables))

	for _, v := range inputVariables {
		dummyInputs[v] = "foo"
	}

	_, err := RenderTemplate(template, templateFormat, dummyInputs)
	return err
}

// RenderTemplate renders the template with the given values.
func RenderTemplate(tmpl string, tmplFormat TemplateFormat, values map[string]any) (string, error) {
	formatter, ok := defaultFormatterMapping[tmplFormat]
	if !ok {
		return "", newInvalidTemplateError(tmplFormat)
	}

	return formatter(tmpl, values)
}
