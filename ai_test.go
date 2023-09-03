package stella_test

import (
	"stella"
	"testing"
)

func TestAI(t *testing.T) {
	ai := stella.New()
	ai.Context.Set("value", 1)

	ai.Func("hello_world", func(args ...any) any {
		v := ai.Context.Get("value").(int)
		ai.Context.Set("value", v+1)
		return v
	})

	err := ai.Prompt("test", "this is a test {{ call .hello_world }} to see if context changes {{ call .hello_world }}")

	if err != nil {
		t.Error(err)
		return
	}

	out, err := ai.Render("test", map[string]any{
		"input": "testing123",
	})

	if err != nil {
		t.Error(err)
		return
	}

	if out != "this is a test 1 to see if context changes 2" {
		t.Errorf("expected 'this is a test 1 to see if context changes 2', received '%s'", out)
		return
	}
}
