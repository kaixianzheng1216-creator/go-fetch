package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func main() {
	outputPath := flag.String("out", "api/openapi.json", "output OpenAPI JSON file")
	flag.Parse()

	openAPIJSON, err := httpapi.OpenAPIJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate OpenAPI: %v\n", err)
		os.Exit(1)
	}
	openAPIJSON = append(openAPIJSON, '\n')

	if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, openAPIJSON, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write OpenAPI file: %v\n", err)
		os.Exit(1)
	}
}
