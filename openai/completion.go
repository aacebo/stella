package openai

type Completion struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
}
