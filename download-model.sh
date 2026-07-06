#!/usr/bin/env sh
set -e

MODEL="${1:-qwen2.5-coder-0.5b}"
OUT="${2:-./models/model.gguf}"

case "$MODEL" in
  qwen2.5-coder-0.5b*)
    URL="https://huggingface.co/Qwen/Qwen2.5-Coder-0.5B-Instruct-GGUF/resolve/main/qwen2.5-coder-0.5b-instruct-q4_k_m.gguf"
    ;;
  qwen2.5-coder-1.5b*)
    URL="https://huggingface.co/Qwen/Qwen2.5-Coder-1.5B-Instruct-GGUF/resolve/main/qwen2.5-coder-1.5b-instruct-q4_k_m.gguf"
    ;;
  qwen2.5-coder-3b*)
    URL="https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf"
    ;;
  llama3.2-1b*)
    URL="https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q4_K_M.gguf"
    ;;
  llama3.2-3b*)
    URL="https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf"
    ;;
  *)
    echo "Unknown model: $MODEL"
    echo "Available: qwen2.5-coder-0.5b, qwen2.5-coder-1.5b, qwen2.5-coder-3b, llama3.2-1b, llama3.2-3b"
    exit 1
    ;;
esac

mkdir -p "$(dirname "$OUT")"
echo "Downloading $MODEL to $OUT..."
curl -L --progress-bar "$URL" -o "$OUT"
echo "Done"
ls -lh "$OUT"
