package boolparser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Parser[any] = Parser{}

// Parser is an output parser used to parse the output of an LLM as a boolean.
type Parser struct {
	TrueStrings  []string
	FalseStrings []string
}

// New returns a bool Parser.
func New() Parser {
	return Parser{
		TrueStrings:  []string{"YES", "TRUE"},
		FalseStrings: []string{"NO", "FALSE"},
	}
}

// FormatInstructions returns instructions on the expected output format.
func (p Parser) FormatInstructions() string {
	return "Your output should be a boolean. e.g.:\n `true` or `false`"
}

func (p Parser) parse(text string) (bool, error) {
	text = normalize(text)

	if slices.Contains(p.TrueStrings, text) {
		return true, nil
	}

	if slices.Contains(p.FalseStrings, text) {
		return false, nil
	}

	return false, llm.ParseError{
		Text: text,
		Reason: fmt.Sprintf("Expected output to one of %v, received %s",
			append(p.TrueStrings, p.FalseStrings...), text),
	}
}

func normalize(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "'\"`")
	text = strings.ToUpper(text)

	return text
}

// Parse parses the output of an LLM into a map of strings.
func (p Parser) Parse(text string) (any, error) {
	return p.parse(text)
}

// ParseWithPrompt does the same as Parse.
func (p Parser) ParseWithPrompt(text string, _ llm.Prompt) (any, error) {
	return p.parse(text)
}

// Type returns the type of the parser.
func (p Parser) Type() string {
	return "boolean_parser"
}
