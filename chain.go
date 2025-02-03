package llm

import "context"

// Chain is the interface all chains must implement.
type Chain interface {
	// Call runs the logic of the chain and returns the output. This method should
	// not be called directly. Use rather the chains.Call, chains.Run or chains.Predict
	// functions that handles the memory and other aspects of the chain.
	Call(ctx context.Context, inputs map[string]any, options ...ChainOption) (map[string]any, error)
	// GetMemory gets the memory of the chain.
	GetMemory() Memory
	// GetInputKeys returns the input keys the chain expects.
	GetInputKeys() []string
	// GetOutputKeys returns the output keys the chain returns.
	GetOutputKeys() []string
}
