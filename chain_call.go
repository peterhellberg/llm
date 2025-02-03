package llm

import (
	"context"
	"fmt"
)

// ChainCall is the standard function used for executing chains.
func ChainCall(ctx context.Context, c Chain, inputValues map[string]any, options ...ChainOption) (map[string]any, error) {
	fullValues := make(map[string]any, 0)

	for key, value := range inputValues {
		fullValues[key] = value
	}

	newValues, err := c.Memory().LoadVariables(ctx, inputValues)
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

	if err = c.Memory().SaveContext(ctx, inputValues, outputValues); err != nil {
		return outputValues, err
	}

	return outputValues, nil
}

func call(ctx context.Context, c Chain, fullValues map[string]any, options ...ChainOption) (map[string]any, error) {
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

func validateInputs(c Chain, inputValues map[string]any) error {
	for _, k := range c.InputKeys() {
		if _, ok := inputValues[k]; !ok {
			return fmt.Errorf("%w: %w: %v", ErrInvalidInputValues, ErrMissingInputValues, k)
		}
	}

	return nil
}

func validateOutputs(c Chain, outputValues map[string]any) error {
	for _, k := range c.OutputKeys() {
		if _, ok := outputValues[k]; !ok {
			return fmt.Errorf("%w: %v", ErrInvalidOutputValues, k)
		}
	}

	return nil
}

func getChainHooks(c Chain) ChainHooks {
	if hh, ok := c.(ChainHooker); ok {
		return hh.ChainHooks()
	}

	return nil
}
