package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/server"
)

func main() {
	outputPath := flag.String("out", "api/openapi.json", "output OpenAPI JSON file")
	flag.Parse()

	openAPIJSON, err := server.OpenAPIJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成 OpenAPI 失败: %v\n", err)
		os.Exit(1)
	}
	openAPIJSON = append(openAPIJSON, '\n')

	if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "创建输出目录失败: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, openAPIJSON, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "写入 OpenAPI 文件失败: %v\n", err)
		os.Exit(1)
	}
}
