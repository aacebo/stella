package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var BASE_URL = "https://api.openai.com"

type Client struct {
	http   http.Client
	apiKey string
}

func NewClient(apiKey string) Client {
	return Client{
		http:   http.Client{},
		apiKey: apiKey,
	}
}

func (self Client) CreateChatCompletion(model string, messages []Message) (any, error) {
	b, err := json.Marshal(map[string]any{
		"model":    model,
		"messages": messages,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", BASE_URL),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.apiKey))
	res, err := self.http.Do(req)

	if err != nil {
		return nil, err
	}

	body := map[string]any{}
	err = json.NewDecoder(res.Body).Decode(&body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
