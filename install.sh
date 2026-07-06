#!/usr/bin/env bash
set -euo pipefail

REPO="jeanmachuca/ai-chat-tui"
BIN="ai-chat"
VERSION="${1:-latest}"

# Detect platform
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

# Determine install dir
if [[ -w /usr/local/bin ]]; then
  DEST="/usr/local/bin"
elif [[ -w ~/.local/bin ]]; then
  DEST="$HOME/.local/bin"
else
  DEST="$HOME/.local/bin"
  mkdir -p "$DEST"
fi

# Prefer building from source if Go is available
if command -v go &>/dev/null; then
  echo "Building from source..."
  TMP="$(mktemp -d)"
  git clone --depth=1 "git@github.com:$REPO.git" "$TMP" 2>/dev/null || \
  git clone --depth=1 "https://github.com/$REPO.git" "$TMP"
  cd "$TMP"
  go build -o "$BIN" .
  cp "$BIN" "$DEST/$BIN"
  rm -rf "$TMP"
  echo "Installed $BIN to $DEST/$BIN (built from source)"
  exit 0
fi

# Fallback: download pre-built release
if [[ "$VERSION" == "latest" ]]; then
  API="https://api.github.com/repos/$REPO/releases/latest"
  URL="$(curl -sL "$API" | grep -oP "browser_download_url.*${OS}_${ARCH}[^\"]*" | cut -d: -f2,3 | tr -d ' "')"
  if [[ -z "$URL" ]]; then
    echo "No release found. Install Go and re-run to build from source."
    exit 1
  fi
else
  URL="https://github.com/$REPO/releases/download/$VERSION/${BIN}_${OS}_${ARCH}.tar.gz"
fi

echo "Downloading $URL..."
curl -fsSL "$URL" -o "/tmp/$BIN.tar.gz"
tar -xzf "/tmp/$BIN.tar.gz" -C /tmp "$BIN"
cp "/tmp/$BIN" "$DEST/$BIN"
rm -f "/tmp/$BIN.tar.gz" "/tmp/$BIN"
echo "Installed $BIN to $DEST/$BIN (release $VERSION)"
