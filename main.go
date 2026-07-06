package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var modelName string
	flag.StringVar(&modelName, "model", "llama3.2", "model name (Ollama) or path (code-inference)")
	var apiBase string
	flag.StringVar(&apiBase, "api", "http://localhost:11434", "API base URL (e.g. http://localhost:11434 for Ollama, http://api:8000 for code-inference)")
	var apiKey string
	flag.StringVar(&apiKey, "api-key", "", "API key for cloud providers (OpenAI, Anthropic, etc.)")
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

	m := initialModel(modelName, apiBase, apiKey, historyFile)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
