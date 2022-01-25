package router

import (
	"io/fs"
	"log"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

// handleFS adds GET handlers for all files and folders using the given filesystem.
func handleFS(r *gin.Engine, fS fs.FS) {
	httpFS := http.FS(fS)
	dirHandler := handleDir(httpFS)

	files, err := fs.ReadDir(fS, ".")
	if err != nil {
		log.Fatal("router.handleFS:", err)
	}

	for _, f := range files {
		name := f.Name()
		if f.IsDir() {
			r.GET("/"+name+"/*"+name, dirHandler)
		} else if name == "index.html" {
			indexHandler := handleIndex(httpFS)
			r.GET("/", indexHandler)
			r.GET("/index.html", indexHandler)
		} else {
			r.GET("/"+name, dirHandler)
		}
	}
}

func handleDir(httpFS http.FileSystem) gin.HandlerFunc {
	fsHandler := http.StripPrefix("/", http.FileServer(httpFS))
	return func(c *gin.Context) {
		fsHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func handleIndex(httpFS http.FileSystem) gin.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		log.Fatal("router.handleIndex:", err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatal("router.handleIndex:", err)
	}

	modtime := stat.ModTime()

	return func(c *gin.Context) {
		http.ServeContent(c.Writer, c.Request, "index.html", modtime, index)
	}
}
