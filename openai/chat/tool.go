package chat

import "encoding/json"

type Tool struct {
	Type     string       `json:"type"`
	Function FunctionTool `json:"function"`
}

type FunctionTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function FunctionToolCall `json:"function"`
}

type FunctionToolCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments,omitempty"`
}

func (self FunctionToolCall) GetName() string {
	return self.Name
}

func (self FunctionToolCall) GetArguments() (any, error) {
	var args any

	if err := json.Unmarshal([]byte(self.Arguments), &args); err != nil {
		return nil, err
	}

	return args, nil
}
