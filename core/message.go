package stella

import "encoding/json"

type Message interface {
	GetRole() string
	GetContent() string
	GetFunctionCalls() []FunctionCall
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatMessage(role string, content string) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: content,
	}
}

func SystemChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "system",
		Content: content,
	}
}

func AssistantChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "assistant",
		Content: content,
	}
}

func UserChatMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "user",
		Content: content,
	}
}

func (self ChatMessage) GetRole() string {
	return self.Role
}

func (self ChatMessage) GetContent() string {
	return self.Content
}

func (self ChatMessage) GetFunctionCalls() []FunctionCall {
	return nil
}

func (self ChatMessage) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
