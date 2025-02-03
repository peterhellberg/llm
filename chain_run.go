package llm

import "context"

// ChainRun can be used to execute a chain if the chain only expects one input and one string output.
func ChainRun(ctx context.Context, c Chain, input any, options ...ChainOption) (string, error) {
	var (
		inputKeys  = c.InputKeys()
		memoryKeys = c.Memory().Variables(ctx)
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
		return "", ErrMultipleInputsInRun
	}

	outputKeys := c.OutputKeys()
	if len(outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{
		neededKeys[0]: input,
	}

	outputValues, err := ChainCall(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return outputValue, nil
}
