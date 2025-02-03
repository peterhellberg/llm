package conversationchain

import (
	"context"
	"testing"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/mock"
)

func TestNew(t *testing.T) {
	var (
		ctx      = context.Background()
		provider = mock.Provider{}
		memory   = mock.Memory{
			VariablesFunc: func(context.Context) []string {
				return []string{"FOO", "BAR"}
			},
		}

		hooks = mock.Hooks{}

		chain = New(provider,
			llm.ChainWithMemory(memory),
			llm.ChainWithHooks(hooks),
		)
	)

	if got, want := len(chain.Memory().Variables(ctx)), 2; got != want {
		t.Fatalf("len(chain.Memory().Variables()) = %d, want %d", got, want)
	}
}
