package gemini

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	stella "github.com/aacebo/stella/core"
	"github.com/aacebo/stella/utils"
)

var chunkFixer = regexp.MustCompile("}\\s*{")

type Completion struct {
	Candidates []CompletionCandidate `json:"candidates"`
}

type CompletionCandidate struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Content      Content `json:"content"`
}

func (self Client) ChatCompletion(messages []stella.Message, stream func(stella.Message)) (stella.Message, error) {
	b, err := json.Marshal(map[string]any{
		"contents": utils.SliceMap(messages, func(message stella.Message) Content {
			role := "user"

			if message.GetRole() != stella.USER && message.GetRole() != stella.SYSTEM {
				role = "model"
			}

			return Content{
				Role:  role,
				Parts: []Part{{Text: message.GetContent()}},
			}
		}),
		"generationConfig": map[string]any{
			"temperature": self.temperature,
		},
	})

	if err != nil {
		return Content{}, err
	}

	method := "generateContent"

	if self.stream {
		method = "streamGenerateContent"
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s:%s?key=%s", BASE_URL, self.model, method, self.apiKey),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return Content{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	if self.stream {
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")
	}

	res, err := self.http.Do(req)

	if err != nil {
		return Content{}, err
	}

	completion := Completion{Candidates: []CompletionCandidate{}}

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
				return Content{}, err
			}

			line = bytes.TrimSpace(line)
			line = bytes.TrimPrefix(line, []byte("data: "))

			if string(line) == "[[DONE]]" {
				break
			}

			chunks := []Completion{}
			err = json.Unmarshal(line, &chunks)

			if err != nil {
				return Content{}, err
			}

			for _, chunk := range chunks {
				if len(chunk.Candidates) == 0 || chunk.Candidates[0].Content.GetContent() == "" {
					continue
				}

				for _, candidate := range chunk.Candidates {
					if stream != nil {
						stream(candidate.Content)
					}

					if candidate.Index > len(completion.Candidates)-1 {
						role := string(stella.ASSISTANT)

						if candidate.Content.Role != "" {
							role = candidate.Content.Role
						}

						completion.Candidates = append(completion.Candidates, CompletionCandidate{
							Content: Content{
								Role:  role,
								Parts: candidate.Content.Parts,
							},
						})
					} else {
						completion.Candidates[candidate.Index].Content.Parts = append(
							completion.Candidates[candidate.Index].Content.Parts,
							candidate.Content.Parts[0],
						)
					}
				}
			}
		}
	} else {
		err = json.NewDecoder(res.Body).Decode(&completion)
	}

	if err != nil {
		return Content{}, err
	}

	if len(completion.Candidates) == 0 {
		return Content{}, errors.New("[gemini.chat] => no message returned")
	}

	return completion.Candidates[0].Content, nil
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
