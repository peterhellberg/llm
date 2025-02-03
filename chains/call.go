package chains

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
)

// Call is the standard function used for executing chains.
func Call(ctx context.Context, c llm.Chain, inputValues map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	fullValues := make(map[string]any, 0)

	for key, value := range inputValues {
		fullValues[key] = value
	}

	newValues, err := c.GetMemory().LoadMemoryVariables(ctx, inputValues)
	if err != nil {
		return nil, err
	}

	for key, value := range newValues {
		fullValues[key] = value
	}

	chainHooks := getChainHooks(c)

	if chainHooks != nil {
		chainHooks.ChainStart(ctx, inputValues)
	}

	outputValues, err := call(ctx, c, fullValues, options...)
	if err != nil {
		if chainHooks != nil {
			chainHooks.ChainError(ctx, err)
		}

		return outputValues, err
	}

	if chainHooks != nil {
		chainHooks.ChainEnd(ctx, outputValues)
	}

	if err = c.GetMemory().SaveContext(ctx, inputValues, outputValues); err != nil {
		return outputValues, err
	}

	return outputValues, nil
}

func call(ctx context.Context, c llm.Chain, fullValues map[string]any, options ...llm.ChainOption) (map[string]any, error) {
	if err := validateInputs(c, fullValues); err != nil {
		return nil, err
	}

	outputValues, err := c.Call(ctx, fullValues, options...)
	if err != nil {
		return outputValues, err
	}

	if err := validateOutputs(c, outputValues); err != nil {
		return outputValues, err
	}

	return outputValues, nil
}

func validateInputs(c llm.Chain, inputValues map[string]any) error {
	for _, k := range c.GetInputKeys() {
		if _, ok := inputValues[k]; !ok {
			return fmt.Errorf("%w: %w: %v", llm.ErrInvalidInputValues, llm.ErrMissingInputValues, k)
		}
	}

	return nil
}

func validateOutputs(c llm.Chain, outputValues map[string]any) error {
	for _, k := range c.GetOutputKeys() {
		if _, ok := outputValues[k]; !ok {
			return fmt.Errorf("%w: %v", llm.ErrInvalidOutputValues, k)
		}
	}

	return nil
}

func getChainHooks(c llm.Chain) llm.ChainHooks {
	if hh, ok := c.(llm.ChainHooker); ok {
		return hh.ChainHooks()
	}

	return nil
}
