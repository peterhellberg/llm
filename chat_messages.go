package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ChatMessageType is the type of chat message.
type ChatMessageType string

const (
	// ChatMessageTypeAI is a message sent by an AI.
	ChatMessageTypeAI ChatMessageType = "ai"
	// ChatMessageTypeHuman is a message sent by a human.
	ChatMessageTypeHuman ChatMessageType = "human"
	// ChatMessageTypeSystem is a message sent by the system.
	ChatMessageTypeSystem ChatMessageType = "system"
	// ChatMessageTypeGeneric is a message sent by a generic user.
	ChatMessageTypeGeneric ChatMessageType = "generic"
	// ChatMessageTypeFunction is a message sent by a function.
	ChatMessageTypeFunction ChatMessageType = "function"
	// ChatMessageTypeTool is a message sent by a tool.
	ChatMessageTypeTool ChatMessageType = "tool"
)

// ChatMessage represents a message in a chat.
type ChatMessage interface {
	// Type gets the type of the message.
	Type() ChatMessageType
	// MessageContent gets the content of the message.
	MessageContent() string
}

// Named is an interface for objects that have a message name.
type Named interface {
	MessageName() string
}

// Statically assert that the types implement the interface.
var (
	_ ChatMessage = AIChatMessage{}
	_ ChatMessage = HumanChatMessage{}
	_ ChatMessage = SystemChatMessage{}
	_ ChatMessage = GenericChatMessage{}
	_ ChatMessage = ToolChatMessage{}
)

// AIChatMessage is a message sent by an AI.
type AIChatMessage struct {
	// Content is the content of the message.
	Content string `json:"content,omitempty"`

	// FunctionCall represents the model choosing to call a function.
	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// ToolCalls represents the model choosing to call tools.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

func (ai AIChatMessage) Type() ChatMessageType              { return ChatMessageTypeAI }
func (ai AIChatMessage) MessageContent() string             { return ai.Content }
func (ai AIChatMessage) MessageFunctionCall() *FunctionCall { return ai.FunctionCall }

// HumanChatMessage is a message sent by a human.
type HumanChatMessage struct {
	Content string
}

func (human HumanChatMessage) Type() ChatMessageType  { return ChatMessageTypeHuman }
func (human HumanChatMessage) MessageContent() string { return human.Content }

// SystemChatMessage is a chat message representing information that should be instructions to the AI system.
type SystemChatMessage struct {
	Content string
}

func (system SystemChatMessage) Type() ChatMessageType  { return ChatMessageTypeSystem }
func (system SystemChatMessage) MessageContent() string { return system.Content }

// GenericChatMessage is a chat message with an arbitrary speaker.
type GenericChatMessage struct {
	Content string
	Role    string
	Name    string
}

func (m GenericChatMessage) Type() ChatMessageType  { return ChatMessageTypeGeneric }
func (m GenericChatMessage) MessageContent() string { return m.Content }
func (m GenericChatMessage) MessageName() string    { return m.Name }

// ToolChatMessage is a chat message representing the result of a tool call.
type ToolChatMessage struct {
	// MessageID is the id of the tool call.
	CallID string `json:"tool_call_id"`
	// Content is the content of the tool message.
	Content string `json:"content"`
}

func (tool ToolChatMessage) Type() ChatMessageType  { return ChatMessageTypeTool }
func (tool ToolChatMessage) MessageContent() string { return tool.Content }
func (tool ToolChatMessage) ID() string             { return tool.CallID }

// GetBufferString gets the buffer string of messages.
func GetBufferString(messages []ChatMessage, humanPrefix string, aiPrefix string) (string, error) {
	result := []string{}

	for _, m := range messages {
		role, err := getMessageRole(m, humanPrefix, aiPrefix)
		if err != nil {
			return "", err
		}

		msg := fmt.Sprintf("%s: %s", role, m.MessageContent())

		if m, ok := m.(AIChatMessage); ok && m.FunctionCall != nil {
			j, err := json.Marshal(m.FunctionCall)
			if err != nil {
				return "", err
			}

			msg = fmt.Sprintf("%s %s", msg, string(j))
		}

		result = append(result, msg)
	}

	return strings.Join(result, "\n"), nil
}

func getMessageRole(m ChatMessage, humanPrefix, aiPrefix string) (string, error) {
	var role string
	switch m.Type() {
	case ChatMessageTypeHuman:
		role = humanPrefix
	case ChatMessageTypeAI:
		role = aiPrefix
	case ChatMessageTypeSystem:
		role = "system"
	case ChatMessageTypeGeneric:
		cgm, ok := m.(GenericChatMessage)
		if !ok {
			return "", fmt.Errorf("%w -%+v", ErrUnexpectedChatMessageType, m)
		}

		role = cgm.Role
	case ChatMessageTypeFunction:
		role = "function"
	case ChatMessageTypeTool:
		role = "tool"
	default:
		return "", ErrUnexpectedChatMessageType
	}

	return role, nil
}

type ChatMessageModelData struct {
	Content string `bson:"content" json:"content"`
	Type    string `bson:"type"    json:"type"`
}

type ChatMessageModel struct {
	Type string               `bson:"type" json:"type"`
	Data ChatMessageModelData `bson:"data" json:"data"`
}

func (c ChatMessageModel) ToChatMessage() ChatMessage {
	switch c.Type {
	case string(ChatMessageTypeAI):
		return AIChatMessage{Content: c.Data.Content}
	case string(ChatMessageTypeHuman):
		return HumanChatMessage{Content: c.Data.Content}
	default:
		return nil
	}
}

// ConvertChatMessageToModel Convert a ChatMessage to a ChatMessageModel.
func ConvertChatMessageToModel(m ChatMessage) ChatMessageModel {
	return ChatMessageModel{
		Type: string(m.Type()),
		Data: ChatMessageModelData{
			Type:    string(m.Type()),
			Content: m.MessageContent(),
		},
	}
}
