package writestream

import (
	"context"
	"fmt"
	"io"

	"github.com/peterhellberg/llm"
)

var _ llm.Hooks = Hooks{}

// Hooks that write stream chunks as strings to the embedded Writer.
type Hooks struct {
	llm.EmptyHooks
	io.Writer
}

func (h Hooks) HandleStreamingFunc(_ context.Context, chunk []byte) {
	fmt.Fprintln(h, string(chunk))
}
