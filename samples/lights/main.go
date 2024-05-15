package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/openai"
)

func main() {
	client := openai.NewClient(
		os.Getenv("OPENAI_API_KEY"),
		"gpt-3.5-turbo",
	).WithTemperature(0).WithStream(true)

	app := stella.New().
		WithLogger(log.Default()).
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
		text := scanner.Text()
		res, err := app.Say("default", text)

		for err != nil {
			res, err = app.Say("default", err.Error())
		}

		if res != "" {
			fmt.Println(res)
		}

		fmt.Print("$: ")
	}
}
