package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	m.viewport.SetContent(renderMessages(m.messages))
	b.WriteString(m.viewport.View())

	b.WriteString("\n")

	if m.loading {
		b.WriteString(fmt.Sprintf("\n  %s Thinking...\n", m.spinner.View()))
	}

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444"))
		b.WriteString(fmt.Sprintf("\n  %s\n\n", errStyle.Render(m.err.Error())))
	}

	b.WriteString("\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n")

	focus := "input"
	if m.focus == focusViewport {
		focus = "viewport"
	}
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
		fmt.Sprintf("  [%s] Ctrl+Enter: send | Tab: focus | Ctrl+N: new | Ctrl+D: delete | q: quit", focus),
	)
	b.WriteString(help)

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}
