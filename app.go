package stella

import (
	"errors"
	"fmt"
	"stella/sync"
)

type Handler func(*Ctx, ...any) error

type App struct {
	ctx        *Ctx
	prompts    *sync.Map[string, *Prompt]
	middleware *sync.Slice[Handler]
}

func New() App {
	return App{
		ctx:        NewCtx(),
		prompts:    sync.NewMap[string, *Prompt](),
		middleware: sync.NewSlice[Handler](),
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

func (self *App) With(middleware ...Handler) *App {
	app := &App{
		ctx:        self.ctx,
		prompts:    self.prompts,
		middleware: self.middleware.Copy(),
	}

	for _, handler := range middleware {
		app.middleware.Push(handler)
	}

	return app
}

func (self *App) Func(name string, method func(*Ctx, ...any) any) *App {
	self.ctx.Set(name, func(args ...any) any {
		for _, handler := range self.middleware.Content() {
			err := handler(self.ctx, args...)

			if err != nil {
				return nil
			}
		}

		return method(self.ctx, args...)
	})

	return self
}

func (self *App) Render(name string, input string) (string, error) {
	prompt := self.prompts.Get(name)

	if prompt == nil {
		return "", errors.New(fmt.Sprintf("prompt \"%s\" not found", name))
	}

	return prompt.Render(self.ctx, input)
}
