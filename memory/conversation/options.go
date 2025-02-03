package conversation

import (
	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/chat"
)

// Option is a function for creating new buffer
// with other than the default values.
type Option func(b *Buffer)

// WithChatHistory is an option for providing the chat history store.
func WithChatHistory(chatHistory llm.ChatMessageHistory) Option {
	return func(b *Buffer) {
		b.chatHistory = chatHistory
	}
}

// WithReturnMessages is an option for specifying should it return messages.
func WithReturnMessages(returnMessages bool) Option {
	return func(b *Buffer) {
		b.returnMessages = returnMessages
	}
}

// WithInputKey is an option for specifying the input key.
func WithInputKey(inputKey string) Option {
	return func(b *Buffer) {
		b.inputKey = inputKey
	}
}

// WithOutputKey is an option for specifying the output key.
func WithOutputKey(outputKey string) Option {
	return func(b *Buffer) {
		b.outputKey = outputKey
	}
}

// WithHumanPrefix is an option for specifying the human prefix.
func WithHumanPrefix(humanPrefix string) Option {
	return func(b *Buffer) {
		b.humanPrefix = humanPrefix
	}
}

// WithAIPrefix is an option for specifying the AI prefix.
func WithAIPrefix(aiPrefix string) Option {
	return func(b *Buffer) {
		b.aiPrefix = aiPrefix
	}
}

// WithMemoryKey is an option for specifying the memory key.
func WithMemoryKey(key string) Option {
	return func(b *Buffer) {
		b.memoryKey = key
	}
}

func applyBufferOptions(opts ...Option) *Buffer {
	m := &Buffer{
		returnMessages: false,
		inputKey:       "",
		outputKey:      "",
		humanPrefix:    "Human",
		aiPrefix:       "AI",
		memoryKey:      "history",
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.chatHistory == nil {
		m.chatHistory = chat.NewMessageHistory()
	}

	return m
}
