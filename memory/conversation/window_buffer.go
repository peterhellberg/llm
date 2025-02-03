package conversation

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.Memory = &WindowBuffer{}

const (
	// defaultConversationWindowSize is the default number of previous conversation.
	defaultConversationWindowSize = 5

	// defaultMessageSize indicates the length of a complete message, currently consisting of 2 parts: ai and human.
	defaultMessageSize = 2
)

// WindowBuffer for storing conversation memory.
type WindowBuffer struct {
	Buffer
	windowSize int
}

// NewWindowBuffer is a function for crating a new window buffer memory.
func NewWindowBuffer(windowSize int, options ...Option) *WindowBuffer {
	if windowSize <= 0 {
		windowSize = defaultConversationWindowSize
	}

	return &WindowBuffer{
		windowSize: windowSize,
		Buffer:     *applyBufferOptions(options...),
	}
}

// Variables uses ConversationBuffer method for memory variables.
func (wb *WindowBuffer) Variables(ctx context.Context) []string {
	return wb.Buffer.Variables(ctx)
}

// LoadVariables uses ConversationBuffer method for loading memory variables.
func (wb *WindowBuffer) LoadVariables(ctx context.Context, _ map[string]any) (map[string]any, error) {
	messages, err := wb.chatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	messages, _ = wb.cutMessages(messages)

	if wb.returnMessages {
		return map[string]any{
			wb.memoryKey: messages,
		}, nil
	}

	bufferString, err := llm.GetBufferString(messages, wb.humanPrefix, wb.aiPrefix)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		wb.memoryKey: bufferString,
	}, nil
}

// SaveContext uses ConversationBuffer method for saving context and prunes memory buffer if needed.
func (wb *WindowBuffer) SaveContext(ctx context.Context, inputValues map[string]any, outputValues map[string]any) error {
	err := wb.Buffer.SaveContext(ctx, inputValues, outputValues)
	if err != nil {
		return err
	}

	messages, err := wb.chatHistory.Messages(ctx)
	if err != nil {
		return err
	}

	if messages, ok := wb.cutMessages(messages); ok {
		err := wb.chatHistory.SetMessages(ctx, messages)
		if err != nil {
			return err
		}
	}

	return nil
}

// Clear uses ConversationBuffer method for clearing buffer memory.
func (wb *WindowBuffer) Clear(ctx context.Context) error {
	return wb.Buffer.Clear(ctx)
}

func (wb *WindowBuffer) cutMessages(message []llm.ChatMessage) ([]llm.ChatMessage, bool) {
	if len(message) > wb.windowSize*defaultMessageSize {
		return message[len(message)-wb.windowSize*defaultMessageSize:], true
	}

	return message, false
}
