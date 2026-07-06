package main

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	userLabel      = lipgloss.NewStyle().Background(lipgloss.Color("#2d5a2d")).Foreground(lipgloss.Color("#ffffff")).Bold(true).Render(" You ")
	assistantLabel = lipgloss.NewStyle().Background(lipgloss.Color("#2d3a6a")).Foreground(lipgloss.Color("#ffffff")).Bold(true).Render(" AI  ")
)

func renderMarkdown(text string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return text
	}
	out, err := r.Render(text)
	if err != nil {
		return text
	}
	return out
}

func renderMessages(msgs []message) string {
	if len(msgs) == 0 {
		return ""
	}
	var b strings.Builder
	for _, msg := range msgs {
		label := userLabel
		if msg.Role == "assistant" {
			label = assistantLabel
		}
		b.WriteString(label)
		b.WriteString("  ")
		b.WriteString(msg.Content)
		b.WriteString("\n\n")
	}
	return b.String()
}
