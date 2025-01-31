package chat

import (
	"context"

	"github.com/peterhellberg/llm"
)

var _ llm.ChatMessageHistory = &MessageHistory{}

// MessageHistory is a struct that stores chat messages.
type MessageHistory struct {
	messages []llm.ChatMessage
}

// NewMessageHistory creates a new MessageHistory using chat message options.
func NewMessageHistory(options ...MessageHistoryOption) *MessageHistory {
	return applyChatOptions(options...)
}

// Messages returns all messages stored.
func (h *MessageHistory) Messages(_ context.Context) ([]llm.ChatMessage, error) {
	return h.messages, nil
}

// AddAIMessage adds an AIMessage to the chat message history.
func (h *MessageHistory) AddAIMessage(_ context.Context, text string) error {
	h.messages = append(h.messages, llm.AIChatMessage{Content: text})
	return nil
}

// AddUserMessage adds a user to the chat message history.
func (h *MessageHistory) AddUserMessage(_ context.Context, text string) error {
	h.messages = append(h.messages, llm.HumanChatMessage{Content: text})
	return nil
}

func (h *MessageHistory) Clear(_ context.Context) error {
	h.messages = make([]llm.ChatMessage, 0)
	return nil
}

func (h *MessageHistory) AddMessage(_ context.Context, message llm.ChatMessage) error {
	h.messages = append(h.messages, message)
	return nil
}

func (h *MessageHistory) SetMessages(_ context.Context, messages []llm.ChatMessage) error {
	h.messages = messages
	return nil
}
