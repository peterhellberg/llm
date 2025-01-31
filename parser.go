package llm

import "fmt"

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
