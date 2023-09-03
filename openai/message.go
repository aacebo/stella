package openai

type MessageRole string

const (
	SYSTEM    MessageRole = "system"
	ASSISTANT MessageRole = "assistant"
	USER      MessageRole = "user"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func NewMessage(role MessageRole, content string) *Message {
	return &Message{
		Role:    role,
		Content: content,
	}
}

func SystemMessage(content string) *Message {
	return &Message{
		Role:    SYSTEM,
		Content: content,
	}
}

func AssistantMessage(content string) *Message {
	return &Message{
		Role:    ASSISTANT,
		Content: content,
	}
}

func UserMessage(content string) *Message {
	return &Message{
		Role:    USER,
		Content: content,
	}
}
