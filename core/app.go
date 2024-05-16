package stella

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Handler func(*Ctx, ...any) error

type App struct {
	ctx        *Ctx
	chat       ChatClient
	prompts    map[string]Prompt
	functions  map[string]Function
	middleware []Handler
	messages   []Message
	logger     *log.Logger
}

func New() *App {
	return &App{
		ctx:        NewCtx(),
		chat:       nil,
		prompts:    map[string]Prompt{},
		functions:  map[string]Function{},
		middleware: []Handler{},
		messages:   []Message{},
		logger:     nil,
	}
}

func (self *App) WithChat(client ChatClient) *App {
	self.chat = client
	return self
}

func (self *App) WithLogger(logger *log.Logger) *App {
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

func (self *App) Func(name string, description string, properties map[string]any, callback FunctionHandler) *App {
	self.functions[name] = Function{
		Name:        name,
		Description: description,
		Properties:  properties,
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

func (self *App) Say(name string, input string, stream func(string)) (string, error) {
	if len(self.messages) == 0 {
		system, err := self.Render(name, input)

		if err != nil {
			return "", err
		}

		self.messages = append(self.messages, SystemChatMessage(system))
	}

	state := map[string]any{}

	for name, def := range self.functions {
		state[name] = def.Handler
	}

	text := ""
	self.messages = append(self.messages, UserChatMessage(input))
	res, err := self.chat.ChatCompletion(self.messages, func(message Message) {
		text += message.GetContent()
		content := text
		prompt, err := NewPrompt("default", content, state)

		if err != nil {
			return
		}

		text = ""
		content, _ = prompt.Render(state)

		if stream != nil {
			stream(content)
		}
	})

	if err != nil {
		return "", err
	}

	prompt, err := NewPrompt("default", res.GetContent(), state)

	if err != nil {
		return "", err
	}

	rendered, err := prompt.Render(state)

	if err != nil {
		return "", err
	}

	self.messages = append(self.messages, NewChatMessage(
		res.GetRole(),
		res.GetContent(),
	))

	if self.logger != nil {
		self.logger.Println(self.messages[len(self.messages)-1])
	}

	return rendered, nil
}
