package stella_test

import (
	"stella"
	"testing"
)

func TestApp(t *testing.T) {
	app := stella.New()
	app.Var("value", 1)

	err := app.Prompt("test", "this is a test {{ call .hello_world }} to see if context changes {{ call .hello_world }}")

	if err != nil {
		t.Error(err)
		return
	}

	out, err := app.With(
		func(ctx *stella.Ctx, args ...any) error {
			v := ctx.Get("value").(int)
			ctx.Set("value", v+1)
			return nil
		},
	).Func("hello_world", func(ctx *stella.Ctx, args ...any) any {
		v := ctx.Get("value").(int)
		ctx.Set("value", v+1)
		return v
	}).Render("test", "testing123")

	if err != nil {
		t.Error(err)
		return
	}

	if out != "this is a test 2 to see if context changes 4" {
		t.Errorf("expected 'this is a test 2 to see if context changes 4', received '%s'", out)
		return
	}
}
