package web

import (
	"embed"
	"io/fs"
)

//go:embed static/* dist
var assets embed.FS

var (
	staticFS = mustSubFS("static")
	distFS   = mustSubFS("dist")
)

func mustSubFS(dir string) fs.FS {
	subFS, err := fs.Sub(assets, dir)
	if err != nil {
		panic("embed " + dir + ": " + err.Error())
	}

	return subFS
}

func StaticFS() fs.FS {
	return staticFS
}

func DistFS() fs.FS {
	return distFS
}

func IndexHTML() ([]byte, error) {
	return assets.ReadFile("dist/index.html")
}
