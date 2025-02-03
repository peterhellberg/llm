package memory

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Memory = Empty{}

// Empty memory that does nothing. Useful for embedding.
type Empty struct{}

func (Empty) Variables(context.Context) []string { return nil }
func (Empty) LoadVariables(context.Context, map[string]any) (map[string]any, error) {
	return make(map[string]any), nil
}
func (Empty) SaveContext(context.Context, map[string]any, map[string]any) error { return nil }
func (Empty) Clear(context.Context) error                                       { return nil }
func (Empty) MemoryKey(context.Context) string                               { return "" }
