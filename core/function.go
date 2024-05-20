package stella

import (
	"fmt"
)

type FunctionHandler func(*Ctx, ...any) (any, error)

type Function struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Properties  map[string]any   `json:"properties"`
	Handler     func(...any) any `json:"-"`
}

func (self Function) String() string {
	return fmt.Sprintf(
		"%s:\n\tdescription: %s\n\treturns: boolean",
		self.Name, self.Description,
	)
}

type FunctionCall interface {
	GetName() string
	GetArguments() (any, error)
}
