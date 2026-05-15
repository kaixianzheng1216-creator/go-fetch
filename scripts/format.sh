#!/usr/bin/env sh
set -eu

root="$(CDPATH= cd "$(dirname "$0")/.." && pwd)"
frontend="$root/frontend"
tracking_script="$root/internal/static/js/script.js"

go_files="$(go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}' ./...)"
if [ -n "$go_files" ]; then
  # shellcheck disable=SC2086
  gofmt -w $go_files
fi

npm --prefix "$frontend" run format
npm --prefix "$frontend" exec prettier -- --write "$tracking_script"
