package conversation

import (
	"context"
	"fmt"

	"github.com/peterhellberg/llm"
)

var _ llm.Memory = &Buffer{}

// Buffer is a simple form of memory that remembers previous conversational back and forth directly.
type Buffer struct {
	chatHistory llm.ChatMessageHistory

	returnMessages bool
	inputKey       string
	outputKey      string
	humanPrefix    string
	aiPrefix       string
	memoryKey      string
}

// NewBuffer is a function for crating a new buffer memory.
func NewBuffer(options ...Option) *Buffer {
	return applyBufferOptions(options...)
}

// Variables gets the input key the buffer memory class will load dynamically.
func (m *Buffer) Variables(context.Context) []string {
	return []string{m.memoryKey}
}

// LoadVariables returns the previous chat messages stored in memory. Previous chat messages
// are returned in a map with the key specified in the MemoryKey field. This key defaults to
// "history". If ReturnMessages is set to true the output is a slice of llms.ChatMessage. Otherwise,
// the output is a buffer string of the chat messages.
func (m *Buffer) LoadVariables(ctx context.Context, _ map[string]any) (map[string]any, error) {
	messages, err := m.chatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	if m.returnMessages {
		return map[string]any{
			m.memoryKey: messages,
		}, nil
	}

	bufferString, err := llm.GetBufferString(messages, m.humanPrefix, m.aiPrefix)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		m.memoryKey: bufferString,
	}, nil
}

// SaveContext uses the input values to the llm to save a user message, and the output values
// of the llm to save an AI message. If the input or output key is not set, the input values or
// output values must contain only one key such that the function can know what string to
// add as a user and AI message. On the other hand, if the output key or input key is set, the
// input key must be a key in the input values and the output key must be a key in the output
// values. The values in the input and output values used to save a user and AI message must
// be strings.
func (m *Buffer) SaveContext(ctx context.Context, inputValues map[string]any, outputValues map[string]any) error {
	userInputValue, err := GetInputValue(inputValues, m.inputKey)
	if err != nil {
		return err
	}

	if err := m.chatHistory.AddUserMessage(ctx, userInputValue); err != nil {
		return err
	}

	aiOutputValue, err := GetInputValue(outputValues, m.outputKey)
	if err != nil {
		return err
	}

	return m.chatHistory.AddAIMessage(ctx, aiOutputValue)
}

// Clear sets the chat messages to a new and empty chat message history.
func (m *Buffer) Clear(ctx context.Context) error {
	return m.chatHistory.Clear(ctx)
}

func (m *Buffer) MemoryKey(context.Context) string {
	return m.memoryKey
}

func GetInputValue(inputValues map[string]any, inputKey string) (string, error) {
	// If the input key is set, return the value in the inputValues with the input key.
	if inputKey != "" {
		inputValue, ok := inputValues[inputKey]
		if !ok {
			return "", fmt.Errorf(
				"%w: %v do not contain inputKey %s",
				llm.ErrInvalidInputValues,
				inputValues,
				inputKey,
			)
		}

		return getInputValueReturnToString(inputValue)
	}

	// Otherwise error if length of map isn't one, or return the only entry in the map.
	if len(inputValues) > 1 {
		return "", fmt.Errorf(
			"%w: multiple keys and no input key set",
			llm.ErrInvalidInputValues,
		)
	}

	for _, inputValue := range inputValues {
		return getInputValueReturnToString(inputValue)
	}

	return "", fmt.Errorf("%w: 0 keys", llm.ErrInvalidInputValues)
}

func getInputValueReturnToString(inputValue any) (string, error) {
	switch value := inputValue.(type) {
	case string:
		return value, nil
	default:
		return "", fmt.Errorf(
			"%w: input value %v not string",
			llm.ErrInvalidInputValues,
			inputValue,
		)
	}
}
