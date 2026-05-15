.PHONY: check format format-check test test-race vet lint frontend-check

GO_FILES := $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)
TRACKING_SCRIPT := internal/static/static/script.js

check: format-check vet lint test frontend-check

format:
	gofmt -w $(GO_FILES)
	npm --prefix frontend run format
	npm --prefix frontend exec prettier -- --write $(TRACKING_SCRIPT)

format-check:
	@files="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$files" ]; then \
		echo "The following Go files need gofmt:"; \
		echo "$$files"; \
		exit 1; \
	fi
	npm --prefix frontend run format:check
	npm --prefix frontend exec prettier -- --check $(TRACKING_SCRIPT)

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

frontend-check:
	npm --prefix frontend run lint
	npm --prefix frontend run build
