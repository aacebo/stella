package stella

type ChatClient interface {
	ChatCompletion(messages []Message, stream func(Message)) (Message, error)
}
