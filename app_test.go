package stella_test

import (
	"log"
	"os"
	"stella"
	"stella/openai"
	"testing"
)

func TestApp(t *testing.T) {
	app := stella.New().WithChat(openai.NewClient(
		os.Getenv("OPENAI_API_KEY"),
		"gpt-3.5-turbo",
	)).Logger(log.Default())

	err := app.Prompt(
		"default",
		"you are an expert on turning the lights on or off and telling me the status.",
	)

	if err != nil {
		t.Error(err)
		return
	}

	app = app.Func("lights_on", "turn the lights on", func(ctx *stella.Ctx, args ...any) (any, error) {
		ctx.Set("state", true)
		return "", nil
	}).Func("lights_off", "turn the lights off", func(ctx *stella.Ctx, args ...any) (any, error) {
		ctx.Set("state", false)
		return "", nil
	}).Func("get_light_status", "get the current light status", func(ctx *stella.Ctx, args ...any) (any, error) {
		return ctx.Get("state", false).(bool), nil
	})

	res, err := app.Say("default", "are the lights on? If not, turn them on and tell me the status.")

	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
