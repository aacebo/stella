package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	stella "github.com/aacebo/stella/core"
)

var chunkFixer = regexp.MustCompile("}\\s*{")

type Completion struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
}

type CompletionChoice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type CompletionChunk struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Model   string                  `json:"model"`
	Choices []CompletionChoiceChunk `json:"choices"`
}

type CompletionChoiceChunk struct {
	Index        int     `json:"index"`
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

func (self Client) CreateChatCompletion(params stella.CreateChatCompletionParams) (stella.Message, error) {
	tools := []Tool{}

	if params.Functions != nil {
		for _, value := range params.Functions {
			tools = append(tools, Tool{
				Type: "function",
				Function: FunctionTool{
					Name:        value.Name,
					Description: value.Description,
					Parameters:  value.Properties,
				},
			})
		}
	}

	b, err := json.Marshal(map[string]any{
		"tool_choice": "auto",
		"model":       self.model,
		"temperature": self.temperature,
		"stream":      self.stream,
		"messages":    params.Messages,
		"tools":       tools,
	})

	if err != nil {
		return Message{}, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", BASE_URL),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.apiKey))

	if self.stream {
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")
	}

	res, err := self.http.Do(req)

	if err != nil {
		return Message{}, err
	}

	completion := Completion{}

	if self.stream {
		scanner := bufio.NewScanner(res.Body)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			// We have a full event payload to parse.
			if i, nlen := containsDoubleNewline(data); i >= 0 {
				return i + nlen, data[0:i], nil
			}
			// If we're at EOF, we have all of the data.
			if atEOF {
				return len(data), data, nil
			}
			// Request more data.
			return 0, nil, nil
		})

		defer res.Body.Close()

		for scanner.Scan() {
			line := scanner.Bytes()

			if err := scanner.Err(); err != nil {
				return Message{}, err
			}

			line = bytes.TrimSpace(line)
			line = bytes.TrimPrefix(line, []byte("data: "))
			line = append([]byte{'['}, append(chunkFixer.ReplaceAll(line, []byte("},{")), ']')...)

			if string(line) == "[[DONE]]" {
				break
			}

			chunks := []CompletionChunk{}
			err = json.Unmarshal(line, &chunks)

			if err != nil {
				return Message{}, err
			}

			for _, chunk := range chunks {
				if len(chunk.Choices) == 0 {
					continue
				}

				completion.ID = chunk.ID
				completion.Object = chunk.Object
				completion.Model = chunk.Model

				for _, choice := range chunk.Choices {
					if params.OnStream != nil && len(choice.Delta.ToolCalls) == 0 {
						params.OnStream(choice.Delta)
					}

					calls := []ToolCall{}

					if choice.Delta.ToolCalls != nil {
						for _, call := range choice.Delta.ToolCalls {
							if call.Valid() {
								calls = append(calls, call)
							}
						}
					}

					if choice.Index > len(completion.Choices)-1 {
						role := "assistant"

						if choice.Delta.Role != "" {
							role = choice.Delta.Role
						}

						completion.Choices = append(completion.Choices, CompletionChoice{
							Message: Message{
								Role:      role,
								Content:   choice.Delta.Content,
								ToolCalls: calls,
							},
						})
					} else {
						completion.Choices[choice.Index].Message.Content += choice.Delta.Content

						if len(calls) > 0 {
							completion.Choices[choice.Index].Message.ToolCalls = calls
						}
					}
				}
			}
		}
	} else {
		err = json.NewDecoder(res.Body).Decode(&completion)
	}

	if err != nil {
		return Message{}, err
	}

	if len(completion.Choices) == 0 {
		return Message{}, errors.New("[openai.chat] => no message returned")
	}

	message := completion.Choices[0].Message

	if message.ToolCalls != nil {
		params.Messages = append(params.Messages, message)

		for _, call := range message.ToolCalls {
			args, err := call.Function.GetArguments()

			if err != nil {
				return message, err
			}

			res := params.Functions[call.Function.Name].Handler(args)
			data, _ := json.Marshal(res)

			params.Messages = append(params.Messages, Message{
				Role:       "tool",
				Content:    string(data),
				ToolCallID: call.ID,
			})

			return self.CreateChatCompletion(params)
		}
	}

	return message, nil
}

func containsDoubleNewline(data []byte) (int, int) {
	// Search for each potentially valid sequence of newline characters
	crcr := bytes.Index(data, []byte("\r\r"))
	lflf := bytes.Index(data, []byte("\n\n"))
	crlflf := bytes.Index(data, []byte("\r\n\n"))
	lfcrlf := bytes.Index(data, []byte("\n\r\n"))
	crlfcrlf := bytes.Index(data, []byte("\r\n\r\n"))
	// Find the earliest position of a double newline combination
	minPos := minPosInt(crcr, minPosInt(lflf, minPosInt(crlflf, minPosInt(lfcrlf, crlfcrlf))))
	// Detemine the length of the sequence
	nlen := 2
	if minPos == crlfcrlf {
		nlen = 4
	} else if minPos == crlflf || minPos == lfcrlf {
		nlen = 3
	}
	return minPos, nlen
}

func minPosInt(a, b int) int {
	if a < 0 {
		return b
	}
	if b < 0 {
		return a
	}
	if a > b {
		return b
	}
	return a
}
