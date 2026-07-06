package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var modelPath string
	flag.StringVar(&modelPath, "model", "/models/model.gguf", "model path on server")
	var apiBase string
	flag.StringVar(&apiBase, "api", "http://api:8000", "API base URL")
	var historyFile string
	flag.StringVar(&historyFile, "history", "", "history file (default: ~/.local/share/ai-chat/history.json)")
	flag.Parse()

	if historyFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		historyFile = home + "/.local/share/ai-chat/history.json"
	}

	m := initialModel(modelPath, apiBase, historyFile)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
