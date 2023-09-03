package stella

import (
	"bytes"
	"text/template"
)

type Prompt struct {
	template *template.Template
}

func NewPrompt(name string, text string) (*Prompt, error) {
	tpl, err := template.New(name).Parse(text)

	if err != nil {
		return nil, err
	}

	return &Prompt{
		template: tpl,
	}, nil
}

func (self Prompt) Render(ctx *Ctx, input string) (string, error) {
	var out bytes.Buffer
	in := ctx.Values()
	in["input"] = input
	err := self.template.Execute(&out, in)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
