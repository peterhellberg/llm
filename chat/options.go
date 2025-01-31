package chat

import "github.com/peterhellberg/llm"

// MessageHistoryOption is a function for creating new chat message history
// with other than the default values.
type MessageHistoryOption func(m *MessageHistory)

// WithPreviousMessages is an option for NewChatMessageHistory for adding
// previous messages to the history.
func WithPreviousMessages(previousMessages []llm.ChatMessage) MessageHistoryOption {
	return func(m *MessageHistory) {
		m.messages = append(m.messages, previousMessages...)
	}
}

func applyChatOptions(options ...MessageHistoryOption) *MessageHistory {
	h := &MessageHistory{
		messages: make([]llm.ChatMessage, 0),
	}

	for _, option := range options {
		option(h)
	}

	return h
}
