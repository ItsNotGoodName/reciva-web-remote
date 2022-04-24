package web

import (
	"embed"
	"io/fs"
	"log"
	"mime"
)

//go:generate npm run build

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

//go:embed dist
var dist embed.FS

func FS() fs.FS {
	f, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal("web.FS:", err)
	}
	return f
}
