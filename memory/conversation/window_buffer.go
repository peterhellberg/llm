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
	WindowSize int
}

// NewWindowBuffer is a function for crating a new window buffer memory.
func NewWindowBuffer(conversationWindowSize int, options ...Option) *WindowBuffer {
	if conversationWindowSize <= 0 {
		conversationWindowSize = defaultConversationWindowSize
	}

	return &WindowBuffer{
		WindowSize: conversationWindowSize,
		Buffer:     *applyBufferOptions(options...),
	}
}

// MemoryVariables uses ConversationBuffer method for memory variables.
func (wb *WindowBuffer) MemoryVariables(ctx context.Context) []string {
	return wb.Buffer.MemoryVariables(ctx)
}

// LoadMemoryVariables uses ConversationBuffer method for loading memory variables.
func (wb *WindowBuffer) LoadMemoryVariables(ctx context.Context, _ map[string]any) (map[string]any, error) {
	messages, err := wb.ChatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	messages, _ = wb.cutMessages(messages)

	if wb.ReturnMessages {
		return map[string]any{
			wb.MemoryKey: messages,
		}, nil
	}

	bufferString, err := llm.GetBufferString(messages, wb.HumanPrefix, wb.AIPrefix)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		wb.MemoryKey: bufferString,
	}, nil
}

// SaveContext uses ConversationBuffer method for saving context and prunes memory buffer if needed.
func (wb *WindowBuffer) SaveContext(ctx context.Context, inputValues map[string]any, outputValues map[string]any) error {
	err := wb.Buffer.SaveContext(ctx, inputValues, outputValues)
	if err != nil {
		return err
	}

	messages, err := wb.ChatHistory.Messages(ctx)
	if err != nil {
		return err
	}

	if messages, ok := wb.cutMessages(messages); ok {
		err := wb.ChatHistory.SetMessages(ctx, messages)
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
	if len(message) > wb.WindowSize*defaultMessageSize {
		return message[len(message)-wb.WindowSize*defaultMessageSize:], true
	}

	return message, false
}
