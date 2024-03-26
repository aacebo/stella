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
	app := stella.New().
		Logger(log.Default()).
		WithChat(openai.NewClient(
			os.Getenv("OPENAI_API_KEY"),
			"gpt-3.5-turbo",
		)).
		Func("lights_on", "turn the lights on", func(ctx *stella.Ctx, args ...any) (any, error) {
			ctx.Set("state", true)
			return "", nil
		}).
		Func("lights_off", "turn the lights off", func(ctx *stella.Ctx, args ...any) (any, error) {
			ctx.Set("state", false)
			return "", nil
		}).
		Func("get_light_status", "get the current light status", func(ctx *stella.Ctx, args ...any) (any, error) {
			return ctx.Get("state", false).(bool), nil
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

		fmt.Println(res)
	}
}
