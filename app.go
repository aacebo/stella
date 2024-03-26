package stella

import (
	"errors"
	"fmt"
	"log"
	"stella/core"
	"strings"
)

type Handler func(*Ctx, ...any) error

type App struct {
	ctx        *Ctx
	chat       ChatClient
	prompts    map[string]Prompt
	functions  map[string]Function
	middleware []Handler
	messages   []core.Message
	logger     *log.Logger
}

func New() *App {
	return &App{
		ctx:        NewCtx(),
		chat:       nil,
		prompts:    map[string]Prompt{},
		functions:  map[string]Function{},
		middleware: []Handler{},
		messages:   []core.Message{},
		logger:     nil,
	}
}

func (self *App) WithChat(client ChatClient) *App {
	self.chat = client
	return self
}

func (self *App) Logger(logger *log.Logger) *App {
	self.logger = logger
	return self
}

func (self *App) Prompt(name string, text string) error {
	functions := map[string]any{}

	for name, def := range self.functions {
		functions[name] = def.Handler
	}

	prompt, err := NewPrompt(name, text, functions)

	if err != nil {
		return err
	}

	self.prompts[name] = prompt
	return nil
}

func (self *App) With(middleware ...Handler) *App {
	for _, handler := range middleware {
		self.middleware = append(self.middleware, handler)
	}

	return self
}

func (self *App) Func(name string, description string, callback FunctionHandler) *App {
	self.functions[name] = Function{
		Name:        name,
		Description: description,
		Handler: func(args ...any) any {
			for _, handler := range self.middleware {
				err := handler(self.ctx, args...)

				if err != nil {
					self.logger.Println(err)
					return ""
				}
			}

			v, err := callback(self.ctx, args...)

			if err != nil {
				self.logger.Println(err)
				return ""
			}

			return v
		},
	}

	return self
}

func (self *App) Render(name string, input string) (string, error) {
	prompt, ok := self.prompts[name]

	if !ok {
		return "", errors.New(fmt.Sprintf("prompt \"%s\" not found", name))
	}

	data := self.ctx.Values()
	data["input"] = input

	for name, def := range self.functions {
		data[name] = def.Handler
	}

	rendered, err := prompt.Render(data)

	if err != nil {
		return "", err
	}

	parts := []string{
		rendered,
		"Do not respond using markdown.",
		"Respond only with template string syntax defined as:",
		"- action syntax examples: https://pkg.go.dev/text/template#hdr-Actions",
		"- variable syntax examples: https://pkg.go.dev/text/template#hdr-Examples",
		"- function syntax examples: https://pkg.go.dev/text/template#hdr-Functions",
	}

	if len(self.ctx.Values()) > 0 {
		parts = append(parts, "Variables:")

		for name := range self.ctx.Values() {
			parts = append(parts, name)
		}
	}

	if len(self.functions) > 0 {
		parts = append(parts, "Functions:")

		for _, function := range self.functions {
			parts = append(parts, function.String())
		}
	}

	return strings.Join(parts, "\n\n"), nil
}

func (self *App) Say(name string, input string) (string, error) {
	if len(self.messages) == 0 {
		system, err := self.Render(name, input)

		if err != nil {
			return "", err
		}

		self.messages = append(self.messages, core.SystemMessage(system))
	}

	self.messages = append(self.messages, core.UserMessage(input))
	res, err := self.chat.ChatCompletion(self.messages)

	if err != nil {
		return "", err
	}

	state := map[string]any{}

	for name, def := range self.functions {
		state[name] = def.Handler
	}

	self.logger.Println(res.Content)
	responsePrompt, err := NewPrompt("default", res.Content, state)

	if err != nil {
		return "", err
	}

	renderedResponse, err := responsePrompt.Render(state)

	if err != nil {
		return "", err
	}

	self.messages = append(self.messages, core.NewMessage(
		res.Role,
		res.Content,
	))

	return renderedResponse, nil
}
