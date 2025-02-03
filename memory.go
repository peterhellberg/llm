package llm

import "context"

// Memory is the interface for memory in chains.
type Memory interface {
	// MemoryKey getter for memory key.
	MemoryKey(ctx context.Context) string
	// Variables Input keys this memory class will load dynamically.
	Variables(ctx context.Context) []string
	// LoadVariables Return key-value pairs given the text input to the chain.
	// If None, return all memories
	LoadVariables(ctx context.Context, inputs map[string]any) (map[string]any, error)
	// SaveContext Save the context of this model run to memory.
	SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error
	// Clear memory contents.
	Clear(ctx context.Context) error
}
