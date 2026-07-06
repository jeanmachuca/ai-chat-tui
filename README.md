# AI Chat TUI

A terminal-based AI chat client built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Works with any OpenAI-compatible API — Ollama, code-inference, OpenAI, or custom endpoints.

## Features

- **Multi-line input** — compose long messages with Shift+Enter for newlines
- **Markdown rendering** — assistant replies rendered with syntax highlighting via [Glamour](https://github.com/charmbracelet/glamour)
- **Session history** — conversations auto-save to `~/.local/share/ai-chat/history.json` (max 50 sessions)
- **SSE streaming** — real-time token-by-token response from compatible APIs
- **Focus toggle** — Tab between input area and scrollable viewport
- **Portable** — same binary works with Ollama, OpenAI, any OpenAI-compatible server

## Usage

```bash
go build -o ai-chat .
```

### Providers

| Provider | Command |
|----------|---------|
| **Ollama** (local) | `./ai-chat --api http://localhost:11434 --model llama3.2` |
| **code-inference** | `./ai-chat --api http://api:8000 --model /models/model.gguf` |
| **OpenAI** | `./ai-chat --api https://api.openai.com/v1 --model gpt-4o-mini --api-key sk-...` |
| **Anthropic** (via proxy) | `./ai-chat --api https://api.anthropic.com/v1 --model claude-sonnet-4-20250514 --api-key sk-ant-...` |
| **Any OpenAI-compatible** | `./ai-chat --api <url> --model <name> [--api-key <key>]` |

### Flags

```
  --api string        API base URL (default "http://localhost:11434")
  --api-key string    API key for cloud providers
  --history string    history file path
  --model string      model name (Ollama) or path (code-inference) (default "llama3.2")
```

## Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+Enter` | Send message |
| `Shift+Enter` | Newline in input |
| `Tab` | Toggle focus: input / viewport |
| `↑` / `↓` | Scroll chat (viewport focused) |
| `Ctrl+N` | New conversation |
| `Ctrl+D` | Delete current conversation |
| `q` / `Ctrl+C` | Quit |

## Architecture

```
main.go       → entry point, flag parsing, program start
model.go      → Bubble Tea model (messages, components, state)
update.go     → event loop (keyboard, SSE stream, session commands)
view.go       → layout (viewport, textarea, spinner, help bar)
api.go        → HTTP client (POST /v1/chat/completions, SSE parsing)
render.go     → Glamour markdown rendering
styles.go     → Lip Gloss color/style definitions
history.go    → JSON session persistence
```

### Data flow

```
User types + Ctrl+Enter
  → update.go appends user message
  → api.go sends POST to /v1/chat/completions (stream: true)
  → SSE response parsed, streamDoneMsg sent to update.go
  → update.go stores assistant message, triggers markdown render
  → view.go renders viewport + textarea + help bar
  → history.go saves session to JSON

Ctrl+N or Ctrl+D
  → update.go triggers saveHistoryCmd / deleteSessionCmd
  → history.go writes updated session store
```

## Requirements

- Go 1.21+
- An OpenAI-compatible API endpoint (Ollama, code-inference, OpenAI, etc.)

No external runtime dependencies — the binary is self-contained.

## Build from source

```bash
git clone git@github.com:jeanmachuca/ai-chat-tui.git
cd ai-chat-tui
go build -o ai-chat .
./ai-chat --api http://localhost:11434 --model llama3.2
```
