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
			fmt.Print("[Gemini]: ")

			res, _ = google.Say("default", text, func(res string) {
				fmt.Print(res)
			})
		} else {
			fmt.Print("[ChatGPT]: ")

			res, _ = openai.Say("default", text, func(res string) {
				fmt.Print(res)
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
