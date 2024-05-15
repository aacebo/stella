package chat

import (
	"encoding/json"

	stella "github.com/aacebo/stella/core"
)

type Message struct {
	Role    stella.MessageRole `json:"role"`
	Content string             `json:"content"`
}

func (self Message) GetRole() stella.MessageRole {
	return self.Role
}

func (self Message) GetContent() string {
	return self.Content
}

func (self Message) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
