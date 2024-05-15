package stella

import "encoding/json"

type MessageRole string

const (
	SYSTEM    MessageRole = "system"
	ASSISTANT MessageRole = "assistant"
	USER      MessageRole = "user"
)

type Message interface {
	GetRole() MessageRole
	GetContent() string
}

type ChatMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func NewChatMessage(role MessageRole, content string) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: content,
	}
}

func SystemChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    SYSTEM,
		Content: content,
	}
}

func AssistantChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    ASSISTANT,
		Content: content,
	}
}

func UserChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    USER,
		Content: content,
	}
}

func (self ChatMessage) GetRole() MessageRole {
	return self.Role
}

func (self ChatMessage) GetContent() string {
	return self.Content
}

func (self ChatMessage) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
