package stella

import (
	"bytes"
	"text/template"
)

type Prompt struct {
	template *template.Template
}

func NewPrompt(name string, text string, functions template.FuncMap) (Prompt, error) {
	tpl, err := template.New(name).Funcs(functions).Parse(text)

	if err != nil {
		return Prompt{}, err
	}

	return Prompt{
		template: tpl,
	}, nil
}

func (self Prompt) Render(in map[string]any) (string, error) {
	var out bytes.Buffer
	err := self.template.Execute(&out, in)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
