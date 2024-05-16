package stella_test

import (
	"log"
	"os"
	"testing"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/openai"
)

func TestApp(t *testing.T) {
	app := stella.New().WithChat(openai.NewChatClient(
		os.Getenv("OPENAI_API_KEY"),
		"gpt-3.5-turbo",
	)).WithLogger(log.Default())

	err := app.Prompt(
		"default",
		"you are an expert on turning the lights on or off and telling me the status.",
	)

	if err != nil {
		t.Error(err)
		return
	}

	app = app.Func("lights_on", "turn the lights on", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
		ctx.Set("state", true)
		return "", nil
	}).Func("lights_off", "turn the lights off", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
		ctx.Set("state", false)
		return "", nil
	}).Func("get_light_status", "get the current light status", nil, func(ctx *stella.Ctx, args ...any) (any, error) {
		return ctx.Get("state", false).(bool), nil
	})

	res, err := app.Say("default", "are the lights on? If not, turn them on and tell me the status.", nil)

	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
