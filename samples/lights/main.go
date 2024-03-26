package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"stella"
	"stella/openai"
)

func main() {
	status := false
	app := stella.New().
		Logger(log.Default()).
		WithChat(openai.NewClient(
			os.Getenv("OPENAI_API_KEY"),
			"gpt-3.5-turbo",
		)).
		Func("lights_on", "turn the lights on", func(ctx *stella.Ctx, args ...any) (any, error) {
			status = true
			return true, nil
		}).
		Func("lights_off", "turn the lights off", func(ctx *stella.Ctx, args ...any) (any, error) {
			status = false
			return false, nil
		}).
		Func("get_light_status", "get the current light status", func(ctx *stella.Ctx, args ...any) (any, error) {
			return status, nil
		})

	err := app.Prompt(
		"default",
		"you are an expert on turning the lights on or off and telling me the status.",
	)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		res, err := app.Say("default", text)

		for err != nil {
			res, err = app.Say("default", err.Error())
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}
