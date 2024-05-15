package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	stella "github.com/aacebo/stella/core"
)

type Completion struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
}

type CompletionChoice struct {
	FinishReason string      `json:"finish_reason"`
	Message      ChatMessage `json:"message"`
}

type CompletionChunk struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Model   string                  `json:"model"`
	Choices []CompletionChoiceChunk `json:"choices"`
}

type CompletionChoiceChunk struct {
	Index int         `json:"index"`
	Delta ChatMessage `json:"delta"`
}

func (self Client) ChatCompletion(messages []stella.Message, stream func(stella.Message)) (stella.Message, error) {
	b, err := json.Marshal(map[string]any{
		"model":       self.model,
		"temperature": self.temperature,
		"stream":      self.stream,
		"messages":    messages,
	})

	if err != nil {
		return ChatMessage{}, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", BASE_URL),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return ChatMessage{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.apiKey))

	if self.stream {
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")
	}

	res, err := self.http.Do(req)

	if err != nil {
		return ChatMessage{}, err
	}

	completion := Completion{}

	if self.stream {
		scanner := bufio.NewScanner(res.Body)
		defer res.Body.Close()

		for scanner.Scan() {
			if scanner.Err(); err != nil {
				if err == io.EOF {
					break
				}

				return ChatMessage{}, err
			}

			data := scanner.Bytes()
			data = bytes.TrimSpace(data)
			data = bytes.TrimPrefix(data, []byte("data: "))

			fmt.Println(string(data))

			if string(data) == "[DONE]" {
				break
			}

			chunk := CompletionChunk{}
			err = json.Unmarshal(data, &chunk)

			if err != nil {
				return ChatMessage{}, err
			}

			completion.ID = chunk.ID
			completion.Object = chunk.Object
			completion.Model = chunk.Model

			for _, choice := range chunk.Choices {
				if stream != nil {
					stream(choice.Delta)
				}

				if choice.Index > len(completion.Choices)-1 {
					completion.Choices = append(completion.Choices, CompletionChoice{
						Message: ChatMessage{
							Role:    stella.MessageRole(choice.Delta.Role),
							Content: choice.Delta.Content,
						},
					})
				} else {
					completion.Choices[choice.Index].Message.Content += choice.Delta.Content
				}
			}
		}
	} else {
		err = json.NewDecoder(res.Body).Decode(&completion)
	}

	if err != nil {
		return ChatMessage{}, err
	}

	if len(completion.Choices) == 0 {
		return ChatMessage{}, errors.New("[openai.chat] => no message returned")
	}

	return completion.Choices[0].Message, nil
}
