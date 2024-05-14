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
		WithLogger(log.Default()).
		WithChat(openai.NewClient(
			os.Getenv("OPENAI_API_KEY"),
			"gpt-4",
		).WithTemperature(0)).
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
