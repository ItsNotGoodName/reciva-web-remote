package web

import (
	"embed"
	"io/fs"
)

//go:generate npm run build

//go:embed dist
var dist embed.FS

func GetFS() fs.FS {
	subFS, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return subFS
}
