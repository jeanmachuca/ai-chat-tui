# AI Chat TUI

Terminal AI chat client — talks to code-inference API at `http://api:8000`.

Built with Bubble Tea, Bubbles, Lip Gloss, Glamour.

## Build & run

```bash
go mod tidy && go build -o ai-chat . && ./ai-chat
```

Flags: `--model`, `--api`, `--history` (see `./ai-chat -h`).

## Keybindings

| Key | Action |
|-----|--------|
| Ctrl+Enter | Send |
| Shift+Enter | Newline |
| Tab | Focus toggle |
| Up/Down | Scroll |
| Ctrl+N | New session |
| Ctrl+D | Delete session |
| q / Ctrl+C | Quit |

## Files

- `main.go` — entry, flags
- `model.go` — Bubble Tea model
- `update.go` — event handling
- `view.go` — layout
- `styles.go` — Lip Gloss
- `api.go` — SSE streaming
- `render.go` — Glamour markdown
- `history.go` — session persistence
