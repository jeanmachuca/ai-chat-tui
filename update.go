package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type streamChunkMsg string

type streamDoneMsg struct {
	fullText string
}

type errMsg struct {
	err error
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 6

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Sequence(saveHistoryCmd(m), tea.Quit)

		case "tab":
			if m.focus == focusInput {
				m.focus = focusViewport
				m.textarea.Blur()
			} else {
				m.focus = focusInput
				m.textarea.Focus()
			}

		case "ctrl+enter":
			if m.loading {
				return m, nil
			}
			text := m.textarea.Value()
			if text == "" {
				return m, nil
			}
			m.messages = append(m.messages, message{Role: "user", Content: text})
			m.streamBuf = ""
			m.err = nil
			m.textarea.Reset()
			m.textarea.Blur()
			m.loading = true

			msgs := make([]message, len(m.messages))
			copy(msgs, m.messages)
			return m, tea.Batch(
				m.spinner.Tick,
				sendMessage(m.apiBase, m.modelPath, msgs),
			)

		case "ctrl+n":
			if !m.loading {
				return m, tea.Sequence(saveHistoryCmd(m), newSessionCmd)
			}

		case "ctrl+d":
			if !m.loading && len(m.messages) > 0 {
				did := m.currentSessID
				return m, tea.Sequence(
					deleteSessionCmd(did, m.historyFile),
					newSessionCmd,
				)
			}
		}

	case streamChunkMsg:
		m.streamBuf += string(msg)
		last := len(m.messages) - 1
		if last >= 0 && m.messages[last].Role == "assistant" {
			m.messages[last].Content = m.streamBuf
		} else {
			m.messages = append(m.messages, message{Role: "assistant", Content: m.streamBuf})
		}
		m.viewport.SetContent(renderMessages(m.messages))
		m.viewport.GotoBottom()

	case streamDoneMsg:
		m.loading = false
		m.focus = focusInput
		m.textarea.Focus()

		if len(m.messages) > 0 {
			last := &m.messages[len(m.messages)-1]
			if last.Role == "assistant" {
				last.Content = renderMarkdown(last.Content)
			}
		}

		m.viewport.SetContent(renderMessages(m.messages))
		m.viewport.GotoBottom()
		return m, saveHistoryCmd(m)

	case errMsg:
		m.loading = false
		m.err = msg.err
		m.textarea.Focus()
		m.focus = focusInput

	case newSessionMsg:
		m.messages = nil
		m.streamBuf = ""
		m.err = nil
		m.loading = false
		m.currentSessID = newID()
		m.textarea.Reset()
		m.textarea.Focus()
		m.focus = focusInput
		m.viewport.SetContent("")
		m.viewport.GotoBottom()
		return m, nil
	}

	if m.focus == focusInput {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var newSessionCmd tea.Cmd = func() tea.Msg {
	return newSessionMsg{}
}

type newSessionMsg struct{}

func saveHistoryCmd(m model) tea.Cmd {
	return func() tea.Msg {
		s := session{
			ID:        m.currentSessID,
			Title:     sessionTitle(m.messages),
			CreatedAt: timeNow(),
			Messages:  m.messages,
		}
		saveSession(m.historyFile, s)
		return nil
	}
}

func deleteSessionCmd(id, path string) tea.Cmd {
	return func() tea.Msg {
		deleteSession(path, id)
		return nil
	}
}
