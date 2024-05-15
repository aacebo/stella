package openai

import "github.com/aacebo/stella/openai/chat"

func NewChatClient(apiKey string, model string) chat.Client {
	return chat.NewClient(apiKey, model)
}
