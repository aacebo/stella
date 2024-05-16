package google

import "github.com/aacebo/stella/google/gemini"

func NewChatClient(apiKey string, model string) gemini.Client {
	return gemini.NewClient(apiKey, model)
}
