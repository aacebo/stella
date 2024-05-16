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
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type CompletionChunk struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Model   string                  `json:"model"`
	Choices []CompletionChoiceChunk `json:"choices"`
}

type CompletionChoiceChunk struct {
	Index int     `json:"index"`
	Delta Message `json:"delta"`
}

func (self Client) ChatCompletion(messages []stella.Message, stream func(stella.Message)) (stella.Message, error) {
	b, err := json.Marshal(map[string]any{
		"model":       self.model,
		"temperature": self.temperature,
		"stream":      self.stream,
		"messages":    messages,
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
				if len(chunk.Choices) == 0 || chunk.Choices[0].Delta.Content == "" {
					continue
				}

				completion.ID = chunk.ID
				completion.Object = chunk.Object
				completion.Model = chunk.Model

				for _, choice := range chunk.Choices {
					if stream != nil {
						stream(choice.Delta)
					}

					if choice.Index > len(completion.Choices)-1 {
						role := stella.ASSISTANT

						if choice.Delta.Role != "" {
							role = choice.Delta.Role
						}

						completion.Choices = append(completion.Choices, CompletionChoice{
							Message: Message{
								Role:    role,
								Content: choice.Delta.Content,
							},
						})
					} else {
						completion.Choices[choice.Index].Message.Content += choice.Delta.Content
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

	return completion.Choices[0].Message, nil
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
