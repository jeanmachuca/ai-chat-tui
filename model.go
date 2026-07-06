package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []message `json:"messages"`
}

type focusMode int

const (
	focusInput focusMode = iota
	focusViewport
)

type model struct {
	messages      []message
	textarea      textarea.Model
	viewport      viewport.Model
	spinner       spinner.Model
	loading       bool
	err           error
	streamBuf     string
	sessions      []session
	currentSessID string
	focus         focusMode
	width         int
	height        int

	modelPath   string
	apiBase     string
	historyFile string
}

func initialModel(modelPath, apiBase, historyFile string) model {
	ta := textarea.New()
	ta.Placeholder = "Type a message (Ctrl+Enter to send)..."
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.Style = viewportStyle

	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Dot

	sessions := loadHistory(historyFile)

	var sessID string
	var msgs []message
	if len(sessions) > 0 {
		latest := sessions[0]
		sessID = latest.ID
		msgs = latest.Messages
	} else {
		sessID = newID()
	}

	return model{
		messages:      msgs,
		textarea:      ta,
		viewport:      vp,
		spinner:       s,
		loading:       false,
		sessions:      sessions,
		currentSessID: sessID,
		focus:         focusInput,
		modelPath:     modelPath,
		apiBase:       apiBase,
		historyFile:   historyFile,
	}
}
