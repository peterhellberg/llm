package chains

import (
	"context"

	"github.com/peterhellberg/llm"
)

// Run can be used to execute a chain if the chain only expects one input and one string output.
func Run(ctx context.Context, c llm.Chain, input any, options ...llm.ChainOption) (string, error) {
	var (
		inputKeys  = c.GetInputKeys()
		memoryKeys = c.GetMemory().MemoryVariables(ctx)
		neededKeys = make([]string, 0, len(inputKeys))
	)

	// Remove keys gotten from the memory.
	for _, inputKey := range inputKeys {
		isInMemory := false

		for _, memoryKey := range memoryKeys {
			if inputKey == memoryKey {
				isInMemory = true
				continue
			}
		}

		if isInMemory {
			continue
		}

		neededKeys = append(neededKeys, inputKey)
	}

	if len(neededKeys) != 1 {
		return "", llm.ErrMultipleInputsInRun
	}

	outputKeys := c.GetOutputKeys()
	if len(outputKeys) != 1 {
		return "", llm.ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{
		neededKeys[0]: input,
	}

	outputValues, err := Call(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", llm.ErrWrongOutputTypeInRun
	}

	return outputValue, nil
}
