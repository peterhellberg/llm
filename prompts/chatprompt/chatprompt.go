package chatprompt

import (
	"fmt"

	"github.com/peterhellberg/llm"
)

var _ llm.Prompt = Value{}

// Value is a prompt value that is a list of chat messages.
type Value []llm.ChatMessage

// String returns the chat message slice as a buffer string.
func (v Value) String() string {
	s, err := llm.BufferString(v, "Human", "AI")
	if err != nil {
		return fmt.Sprintf("%v", []llm.ChatMessage(v))
	}

	return s
}

// Messages returns the ChatMessage slice.
func (v Value) Messages() []llm.ChatMessage {
	return v
}
