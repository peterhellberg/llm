package llm

import (
	"context"
)

// ChainPredict can be used to execute a chain if the chain only expects one string output.
func ChainPredict(ctx context.Context, c Chain, inputValues map[string]any, options ...ChainOption) (string, error) {
	outputValues, err := ChainCall(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputKeys := c.OutputKeys()
	if len(outputKeys) != 1 {
		return "", ErrMultipleOutputsInPredict
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", ErrOutputNotStringInPredict
	}

	return outputValue, nil
}
