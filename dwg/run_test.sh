#!/usr/bin/env bash
# Runs dwg2png against dwg/test.dwg and writes a PNG preview next to it.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

IN="test.dwg"
OUT="test_preview.png"

go run ./cmd/dwg2png "$IN" "$OUT"

echo "wrote $(pwd)/$OUT"
