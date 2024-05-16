package main

import (
	"os"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/google"
)

func Google() *stella.App {
	client := google.NewChatClient(
		os.Getenv("GEMINI_API_KEY"),
		"gemini-1.5-flash-latest",
	).WithTemperature(0).WithStream(STREAM)

	app := stella.New().WithChat(client)
	app.Prompt("default", "you are a chatty robot")

	return app
}
