//go:build !dev

package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed dist
var dist embed.FS

func FS() fs.FS {
	f, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatalln("web.FS:", err)
	}
	return f
}
