package llm

import (
	"fmt"
	"strings"
)

// Parser is an interface for parsing the output of an LLM call.
type Parser[T any] interface {
	// Parse parses the output of an LLM call.
	Parse(text string) (T, error)
	// ParseWithPrompt parses the output of an LLM call with the prompt used.
	ParseWithPrompt(text string, prompt Prompt) (T, error)
	// FormatInstructions returns a string describing the format of the output.
	FormatInstructions() string
	// Type returns the string type key uniquely identifying this class of parser
	Type() string
}

// ParseError is the error type returned by output parsers.
type ParseError struct {
	Text   string
	Reason string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse text %s. %s", e.Text, e.Reason)
}

var _ Parser[any] = EmptyParser{}

// EmptyParser that does nothing. Useful for embedding.
type EmptyParser struct{}

func (EmptyParser) FormatInstructions() string {
	return ""
}

func (EmptyParser) Parse(text string) (any, error) {
	return strings.TrimSpace(text), nil
}

func (EmptyParser) ParseWithPrompt(text string, _ Prompt) (any, error) {
	return strings.TrimSpace(text), nil
}

func (EmptyParser) Type() string {
	return "empty_parser"
}
