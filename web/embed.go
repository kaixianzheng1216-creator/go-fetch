package web

import (
	"embed"
	"io/fs"
)

//go:embed dist static/tracker.js
var assets embed.FS

var dashboardFS = mustSubFS("dist")

func mustSubFS(dir string) fs.FS {
	subFS, err := fs.Sub(assets, dir)
	if err != nil {
		panic("embed assets " + dir + ": " + err.Error())
	}
	return subFS
}

func DashboardFS() fs.FS {
	return dashboardFS
}

func IndexHTML() ([]byte, error) {
	return assets.ReadFile("dist/index.html")
}

func TrackerScript() ([]byte, error) {
	return assets.ReadFile("static/tracker.js")
}
