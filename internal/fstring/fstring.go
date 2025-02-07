package fstring

import "fmt"

var (
	ErrEmptyExpression       = fmt.Errorf("empty expression not allowed")
	ErrArgsNotDefined        = fmt.Errorf("args not defined")
	ErrLeftBracketNotClosed  = fmt.Errorf("single '{' is not allowed")
	ErrRightBracketNotClosed = fmt.Errorf("single '}' is not allowed")
)

// Format interpolates the given template with the given values by using f-string.
func Format(template string, values map[string]any) (string, error) {
	p := newParser(template, values)

	if err := p.parse(); err != nil {
		return "", err
	}

	return string(p.result), nil
}
