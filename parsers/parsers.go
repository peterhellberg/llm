package parsers

import (
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Parser[any] = Empty{}

// Empty parser that does nothing. Useful for embedding.
type Empty struct{}

func (p Empty) FormatInstructions() string {
	return ""
}

func (p Empty) Parse(text string) (any, error) {
	return strings.TrimSpace(text), nil
}

func (p Empty) ParseWithPrompt(text string, _ llm.Prompt) (any, error) {
	return strings.TrimSpace(text), nil
}

func (p Empty) Type() string {
	return "void_parser"
}
