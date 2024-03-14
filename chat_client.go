package stella

import "stella/core"

type ChatClient interface {
	ChatCompletion(messages []core.Message) (core.Message, error)
}
