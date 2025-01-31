package llm

// Message is the content of a message sent to a LLM. It has a role and a
// sequence of parts. For example, it can represent one message in a chat
// session sent by the user, in which case Role will be
// ChatMessageTypeHuman and Parts will be the sequence of items sent in
// this specific message.
type Message struct {
	Role  ChatMessageType
	Parts []ContentPart
}

// TextPartsMessage is a helper function to create a Message with a role and a
// list of text parts.
func TextPartsMessage(role ChatMessageType, parts ...string) Message {
	result := Message{
		Role:  role,
		Parts: []ContentPart{},
	}

	for _, part := range parts {
		result.Parts = append(result.Parts, TextPart(part))
	}

	return result
}
