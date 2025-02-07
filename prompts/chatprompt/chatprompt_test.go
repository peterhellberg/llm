package chatprompt

import (
	"testing"

	"github.com/peterhellberg/llm"
)

func TestValueString(t *testing.T) {
	v := Value([]llm.ChatMessage{
		llm.AIChatMessage{Content: "🤖"},
		llm.HumanChatMessage{Content: "🧍"},
	})

	if got, want := v.String(), "AI: 🤖\nHuman: 🧍"; got != want {
		t.Fatalf("v.String() = %q, want %q", got, want)
	}
}
