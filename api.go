package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type chatRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	Stream    bool      `json:"stream"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type streamLine struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func sendMessage(apiBase, apiKey, modelName string, msgs []message) tea.Cmd {
	return func() tea.Msg {
		body := chatRequest{
			Model:     modelName,
			Messages:  msgs,
			Stream:    true,
			MaxTokens: 4096,
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return errMsg{err: fmt.Errorf("encode: %w", err)}
		}

		req, err := http.NewRequest("POST", apiBase+"/v1/chat/completions", &buf)
		if err != nil {
			return errMsg{err: fmt.Errorf("request: %w", err)}
		}
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errMsg{err: fmt.Errorf("request: %w", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			return errMsg{err: fmt.Errorf("API %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))}
		}

		var fullText strings.Builder
		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
			var sl streamLine
			if err := json.Unmarshal([]byte(data), &sl); err != nil {
				continue
			}
			for _, ch := range sl.Choices {
				if ch.Delta.Content != "" {
					fullText.WriteString(ch.Delta.Content)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return errMsg{err: fmt.Errorf("stream: %w", err)}
		}

		return streamDoneMsg{fullText: fullText.String()}
	}
}
