package main

import "github.com/charmbracelet/lipgloss"

var (
	viewportStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#3a3a3a"))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64b5f6"))
)
