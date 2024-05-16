package main

import (
	"fmt"
)

var STREAM = true

func main() {
	text := "hi"
	turn := 0
	openai := OpenAI()
	google := Google()

	for {
		var res string

		if turn%2 == 0 {
			res, _ = google.Say("default", text, func(res string) {
				fmt.Printf("[Gemini]: %s", res)
			})
		} else {
			res, _ = openai.Say("default", text, func(res string) {
				fmt.Printf("[ChatGPT]: %s", res)
			})
		}

		text = res
		turn++

		if !STREAM {
			fmt.Print(res)
		}

		fmt.Print("\n")
	}
}
