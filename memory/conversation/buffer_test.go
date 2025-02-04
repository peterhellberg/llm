package conversation

import (
	"context"
	"testing"

	"github.com/peterhellberg/llm/chat"
)

func TestNewBuffer(t *testing.T) {
	ctx := context.Background()

	mh := chat.NewMessageHistory()

	ik := "I"
	mk := "M"
	ok := "O"

	b := NewBuffer(
		WithAIPrefix("A"),
		WithChatHistory(mh),
		WithHumanPrefix("H"),
		WithInputKey(ik),
		WithMemoryKey(mk),
		WithOutputKey(ok),
		WithReturnMessages(false),
	)

	if got, want := b.MemoryKey(ctx), mk; got != want {
		t.Fatalf("b.MemoryKey(ctx) = %q, want %q", got, want)
	}

	b.Clear(ctx)

	{
		vars, err := b.LoadVariables(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(vars) != 1 {
			t.Fatalf("unexpected vars: %v", vars)
		}
	}

	{
		if err := b.SaveContext(ctx, map[string]any{
			"I": "FOO",
		}, map[string]any{
			"O": "BAR",
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	{
		vars, err := b.LoadVariables(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(vars) != 1 {
			t.Fatalf("unexpected vars: %v", vars)
		}
	}
}
