package chat

import (
	"encoding/json"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/utils"
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

func (self Message) GetRole() string {
	return self.Role
}

func (self Message) GetContent() string {
	return self.Content
}

func (self Message) GetFunctionCalls() []stella.FunctionCall {
	return utils.SliceMap(self.ToolCalls, func(call ToolCall) stella.FunctionCall {
		return call.Function
	})
}

func (self Message) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
