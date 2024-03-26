package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stella/core"
)

var BASE_URL = "https://api.openai.com"

type Client struct {
	http   http.Client
	apiKey string
	model  string
}

func NewClient(apiKey string, model string) Client {
	return Client{
		http:   http.Client{},
		apiKey: apiKey,
		model:  model,
	}
}

func (self Client) ChatCompletion(messages []core.Message) (core.Message, error) {
	b, err := json.Marshal(map[string]any{
		"model":       self.model,
		"temperature": 0.8,
		"messages":    messages,
	})

	if err != nil {
		return core.Message{}, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", BASE_URL),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return core.Message{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.apiKey))
	res, err := self.http.Do(req)

	if err != nil {
		return core.Message{}, err
	}

	body := Completion{}
	err = json.NewDecoder(res.Body).Decode(&body)

	if err != nil {
		return core.Message{}, err
	}

	if len(body.Choices) == 0 {
		return core.Message{}, errors.New("[openai.chat] => no message returned")
	}

	return body.Choices[0].Message, nil
}
