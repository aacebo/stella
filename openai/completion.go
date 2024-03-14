package openai

import "stella/core"

type Completion struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
}

type CompletionChoice struct {
	FinishReason string       `json:"finish_reason"`
	Message      core.Message `json:"message"`
}
