package stella

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type AI struct {
	Context      Map[string, any]
	conversation Slice[*Message]
	prompts      Map[string, *template.Template]
	functions    Map[string, func(...any) any]
}

func New() AI {
	return AI{
		Context:      NewMap[string, any](),
		conversation: NewSlice[*Message](),
		prompts:      NewMap[string, *template.Template](),
		functions:    NewMap[string, func(...any) any](),
	}
}

func (self *AI) Prompt(name string, content string) error {
	tpl, err := template.New(name).Parse(content)

	if err != nil {
		return err
	}

	self.prompts.Set(name, tpl)
	return nil
}

func (self *AI) Func(name string, method func(...any) any) {
	self.functions.Set(name, method)
}

func (self *AI) Render(name string, ctx map[string]any) (string, error) {
	tpl := self.prompts.Get(name)

	if tpl == nil {
		return "", errors.New(fmt.Sprintf("prompt \"%s\" not found", name))
	}

	in := self.Context.Map()

	for k, v := range self.functions.Map() {
		if _, ok := in[k]; ok {
			return "", errors.New(fmt.Sprintf("duplicate context key \"%s\"", k))
		}

		in[k] = v
	}

	for k, v := range ctx {
		if _, ok := in[k]; ok {
			return "", errors.New(fmt.Sprintf("duplicate context key \"%s\"", k))
		}

		in[k] = v
	}

	var out bytes.Buffer
	err := tpl.Execute(&out, in)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
