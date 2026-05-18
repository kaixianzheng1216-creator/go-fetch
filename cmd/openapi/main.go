package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	outputPath := flag.String("out", "api/openapi.json", "output OpenAPI JSON file")
	flag.Parse()

	openAPIJSON, err := httpapi.OpenAPIJSON()
	if err != nil {
		return fmt.Errorf("generate OpenAPI: %w", err)
	}
	openAPIJSON = append(openAPIJSON, '\n')

	if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(*outputPath, openAPIJSON, 0o644); err != nil {
		return fmt.Errorf("write OpenAPI file: %w", err)
	}

	return nil
}
