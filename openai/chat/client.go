package chat

import (
	"net/http"
)

var BASE_URL = "https://api.openai.com"

type Client struct {
	http        http.Client
	apiKey      string
	model       string
	temperature float32
	stream      bool
}

func NewClient(apiKey string, model string) Client {
	return Client{
		http:        http.Client{},
		apiKey:      apiKey,
		model:       model,
		temperature: 0.8,
		stream:      false,
	}
}

func (self Client) WithTemperature(temperature float32) Client {
	self.temperature = temperature
	return self
}

func (self Client) WithStream(stream bool) Client {
	self.stream = stream
	return self
}

func (self Client) WithHttpClient(client http.Client) Client {
	self.http = client
	return self
}
