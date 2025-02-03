package regexparser

import (
	"fmt"
	"regexp"

	"github.com/peterhellberg/llm"
)

var _ llm.Parser[any] = Parser{}

// New returns a regexp Parser.
func New(str string) Parser {
	expression := regexp.MustCompile(str)
	outputKeys := expression.SubexpNames()[1:]

	return Parser{
		expression: expression,
		outputKeys: outputKeys,
	}
}

// Parser is an output parser used to parse the output of an LLM as a map.
type Parser struct {
	expression *regexp.Regexp
	outputKeys []string
}

// FormatInstructions returns instructions on the expected output format.
func (p Parser) FormatInstructions() string {
	instructions := "Your output should be a map of strings. e.g.:\n"
	instructions += "map[string]string{\"key1\": \"value1\", \"key2\": \"value2\"}"

	return instructions
}

// Type returns the type of the parser.
func (p Parser) Type() string {
	return "regex_parser"
}

// Parse parses the output of an LLM into a map of strings.
func (p Parser) Parse(text string) (any, error) {
	return p.parse(text)
}

// ParseWithPrompt does the same as Parse.
func (p Parser) ParseWithPrompt(text string, _ llm.Prompt) (any, error) {
	return p.parse(text)
}

func (p Parser) parse(text string) (map[string]string, error) {
	match := p.expression.FindStringSubmatch(text)

	if len(match) == 0 {
		return nil, llm.ParseError{
			Text:   text,
			Reason: fmt.Sprintf("No match found for expression %s", p.expression),
		}
	}

	// remove the first match (entire string) for parity with the output keys.
	match = match[1:]

	matches := make(map[string]string, len(match))

	for i, name := range p.outputKeys {
		matches[name] = match[i]
	}

	return matches, nil
}
