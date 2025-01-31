package llm

import "fmt"

func ResolvePartialValues(partialValues map[string]any, values map[string]any) (map[string]any, error) {
	resolvedValues := make(map[string]any)

	for variable, value := range partialValues {
		switch value := value.(type) {
		case string:
			resolvedValues[variable] = value
		case func() string:
			resolvedValues[variable] = value()
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidPartialVariableType, variable)
		}
	}

	for variable, value := range values {
		resolvedValues[variable] = value
	}

	return resolvedValues, nil
}
