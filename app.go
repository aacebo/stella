package stella

import (
	"errors"
	"fmt"
	"stella/sync"
)

type App struct {
	ctx     *Ctx
	prompts sync.Map[string, *Prompt]
}

func New() App {
	return App{
		ctx:     NewCtx(),
		prompts: sync.NewMap[string, *Prompt](),
	}
}

func (self *App) Prompt(name string, text string) error {
	prompt, err := NewPrompt(name, text)

	if err != nil {
		return err
	}

	self.prompts.Set(name, prompt)
	return nil
}

func (self *App) Var(name string, value any) {
	self.ctx.Set(name, value)
}

func (self *App) Func(name string, method func(*Ctx, ...any) any) {
	self.ctx.Set(name, func(args ...any) any {
		return method(self.ctx, args...)
	})
}

func (self *App) Render(name string, input string) (string, error) {
	prompt := self.prompts.Get(name)

	if prompt == nil {
		return "", errors.New(fmt.Sprintf("prompt \"%s\" not found", name))
	}

	return prompt.Render(self.ctx, input)
}
