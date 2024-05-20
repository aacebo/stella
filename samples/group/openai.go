package main

import (
	"os"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/openai"
)

func OpenAI() *stella.App {
	client := openai.NewChatClient(
		os.Getenv("OPENAI_API_KEY"),
		"gpt-4-turbo",
	).WithTemperature(0).WithStream(STREAM)

	app := stella.New().WithChat(client)
	app.Prompt("default", "you are a chatty robot")

	return app
}
