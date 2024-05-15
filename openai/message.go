package openai

import (
	"encoding/json"

	stella "github.com/aacebo/stella/core"
)

type ChatMessage struct {
	Role    stella.MessageRole `json:"role"`
	Content string             `json:"content"`
}

func NewChatMessage(role stella.MessageRole, content string) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: content,
	}
}

func (self ChatMessage) GetRole() stella.MessageRole {
	return self.Role
}

func (self ChatMessage) GetContent() string {
	return self.Content
}

func (self ChatMessage) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
