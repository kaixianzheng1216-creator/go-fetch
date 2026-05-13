package web

import (
	"embed"
	"io/fs"
)

//go:embed static/* dist
var assets embed.FS

var (
	staticFS, _ = fs.Sub(assets, "static")
	distFS, _   = fs.Sub(assets, "dist")
)

func StaticFS() fs.FS {
	return staticFS
}

func DistFS() fs.FS {
	return distFS
}

func IndexHTML() ([]byte, error) {
	return assets.ReadFile("dist/index.html")
}
