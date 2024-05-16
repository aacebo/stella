package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/google"
)

var STREAM = true

func main() {
	// client := openai.NewChatClient(
	// 	os.Getenv("OPENAI_API_KEY"),
	// 	"gpt-4-turbo",
	// ).WithTemperature(0).WithStream(true)
	client := google.NewChatClient(
		os.Getenv("GEMINI_API_KEY"),
		"gemini-1.5-flash-latest",
	).WithTemperature(0).WithStream(STREAM)

	app := stella.New().
		WithChat(client).
		Func("lights_on", "turn the lights on", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
			status := ctx.Get("status", false)

			if status == false {
				status = true
				ctx.Set("status", true)
			}

			return true, nil
		}).
		Func("lights_off", "turn the lights off", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
			status := ctx.Get("status", false)

			if status == true {
				status = false
				ctx.Set("status", false)
			}

			return true, nil
		}).
		Func("get_light_status", "get the current light status", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
			return ctx.Get("status", false), nil
		})

	err := app.Prompt(
		"default",
		"you are an expert on turning the lights on or off and telling me the status.",
	)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("$: ")

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "exit" {
			return
		}

		res, err := app.Say("default", text, func(text string) {
			fmt.Print(text)
		})

		for err != nil {
			fmt.Println(err.Error())
			res, err = app.Say("default", err.Error(), func(text string) {
				fmt.Print(text)
			})
		}

		if !STREAM {
			fmt.Print(res)
		}

		fmt.Print("\n$: ")
	}
}
