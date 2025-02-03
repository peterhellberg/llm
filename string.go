package llm

var _ Prompt = String("")

// String is a prompt value that is a string.
type String string

func (s String) String() string {
	return string(s)
}

// Messages returns a single-element ChatMessage slice.
func (s String) Messages() []ChatMessage {
	return []ChatMessage{
		HumanChatMessage{
			Content: string(s),
		},
	}
}
