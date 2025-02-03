package boolparser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Parser[any] = Parser{}

const (
	parserType               = "bool_parser"
	parserFormatInstructions = "Your output should be a boolean. e.g.:\n `true` or `false`"
)

// Parser is an output parser used to parse the output of an LLM as a boolean.
type Parser struct {
	trueStrings  []string
	falseStrings []string
}

// New returns a bool Parser.
func New() Parser {
	return Parser{
		trueStrings:  []string{"YES", "TRUE"},
		falseStrings: []string{"NO", "FALSE"},
	}
}

// Type returns the type of the parser.
func (p Parser) Type() string {
	return parserType
}

// FormatInstructions returns instructions on the expected output format.
func (p Parser) FormatInstructions() string {
	return parserFormatInstructions
}

// Parse parses the output of an LLM into a map of strings.
func (p Parser) Parse(text string) (any, error) {
	return p.parse(text)
}

func (p Parser) parse(text string) (bool, error) {
	text = normalize(text)

	if slices.Contains(p.trueStrings, text) {
		return true, nil
	}

	if slices.Contains(p.falseStrings, text) {
		return false, nil
	}

	return false, llm.ParseError{
		Text: text,
		Reason: fmt.Sprintf("Expected output to one of %v, received %s",
			append(p.trueStrings, p.falseStrings...), text),
	}
}

func normalize(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "'\"`")
	text = strings.ToUpper(text)

	return text
}
