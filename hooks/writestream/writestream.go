package writestream

import (
	"context"
	"fmt"
	"io"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/hooks"
)

var _ llm.Hooks = Hooks{}

// Hooks that write stream chunks as strings to the embedded Writer.
type Hooks struct {
	hooks.Empty
	io.Writer
}

func (h Hooks) HandleStreamingFunc(_ context.Context, chunk []byte) {
	fmt.Fprintln(h, string(chunk))
}
