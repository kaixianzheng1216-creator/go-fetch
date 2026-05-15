package static

import (
	"embed"
	"io/fs"
)

//go:embed js/* dist
var assets embed.FS

var (
	staticFS = mustSubFS("js")
	distFS   = mustSubFS("dist")
)

func mustSubFS(dir string) fs.FS {
	subFS, err := fs.Sub(assets, dir)
	if err != nil {
		panic("嵌入资源 " + dir + " 失败: " + err.Error())
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
