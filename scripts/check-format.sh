#!/usr/bin/env sh
set -eu

root="$(CDPATH= cd "$(dirname "$0")/.." && pwd)"
dashboard="$root/web/dashboard"
tracking_script="$root/web/static/tracker.js"

go_files="$(go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}' ./...)"
if [ -n "$go_files" ]; then
  # shellcheck disable=SC2086
  unformatted="$(gofmt -l $go_files)"
  if [ -n "$unformatted" ]; then
    printf 'The following Go files need gofmt:\n%s\n' "$unformatted" >&2
    exit 1
  fi
fi

npm --prefix "$dashboard" run format:check
npm --prefix "$dashboard" exec prettier -- --check "$tracking_script"
