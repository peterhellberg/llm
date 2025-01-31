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
		b.ChatHistory = chatHistory
	}
}

// WithReturnMessages is an option for specifying should it return messages.
func WithReturnMessages(returnMessages bool) Option {
	return func(b *Buffer) {
		b.ReturnMessages = returnMessages
	}
}

// WithInputKey is an option for specifying the input key.
func WithInputKey(inputKey string) Option {
	return func(b *Buffer) {
		b.InputKey = inputKey
	}
}

// WithOutputKey is an option for specifying the output key.
func WithOutputKey(outputKey string) Option {
	return func(b *Buffer) {
		b.OutputKey = outputKey
	}
}

// WithHumanPrefix is an option for specifying the human prefix.
func WithHumanPrefix(humanPrefix string) Option {
	return func(b *Buffer) {
		b.HumanPrefix = humanPrefix
	}
}

// WithAIPrefix is an option for specifying the AI prefix.
func WithAIPrefix(aiPrefix string) Option {
	return func(b *Buffer) {
		b.AIPrefix = aiPrefix
	}
}

// WithMemoryKey is an option for specifying the memory key.
func WithMemoryKey(memoryKey string) Option {
	return func(b *Buffer) {
		b.MemoryKey = memoryKey
	}
}

func applyBufferOptions(opts ...Option) *Buffer {
	m := &Buffer{
		ReturnMessages: false,
		InputKey:       "",
		OutputKey:      "",
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.ChatHistory == nil {
		m.ChatHistory = chat.NewMessageHistory()
	}

	return m
}
