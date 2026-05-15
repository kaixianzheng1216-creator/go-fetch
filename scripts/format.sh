#!/usr/bin/env sh
set -eu

root="$(CDPATH= cd "$(dirname "$0")/.." && pwd)"
dashboard="$root/web/dashboard"
tracking_script="$root/web/tracker/script.js"

go_files="$(go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}' ./...)"
if [ -n "$go_files" ]; then
  # shellcheck disable=SC2086
  gofmt -w $go_files
fi

npm --prefix "$dashboard" run format
npm --prefix "$dashboard" exec prettier -- --write "$tracking_script"
