package gemini

import stella "github.com/aacebo/stella/core"

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

func (self Content) GetRole() string {
	return self.Role
}

func (self Content) GetContent() string {
	content := ""

	for _, part := range self.Parts {
		content += part.Text
	}

	return content
}

func (self Content) GetFunctionCalls() []stella.FunctionCall {
	return nil
}

type Part struct {
	Text       string `json:"text,omitempty"`
	InlineData []byte `json:"inlineData,omitempty"`
}
