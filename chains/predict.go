package chains

import (
	"context"

	"github.com/peterhellberg/llm"
)

// Predict can be used to execute a chain if the chain only expects one string output.
func Predict(ctx context.Context, c llm.Chain, inputValues map[string]any, options ...llm.ChainOption) (string, error) {
	outputValues, err := Call(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputKeys := c.OutputKeys()
	if len(outputKeys) != 1 {
		return "", llm.ErrMultipleOutputsInPredict
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", llm.ErrOutputNotStringInPredict
	}

	return outputValue, nil
}
