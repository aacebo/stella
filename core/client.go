package stella

type CreateChatCompletionParams struct {
	Messages  []Message           `json:"messages"`
	Functions map[string]Function `json:"functions"`
	OnStream  func(Message)       `json:"-"`
}

type ChatClient interface {
	SupportsNativeFunctions() bool
	CreateChatCompletion(params CreateChatCompletionParams) (Message, error)
}
