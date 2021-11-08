package routes

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func newWebRouteHandler(r *gin.Engine, httpFS http.FileSystem, httpPath, fsPath string) {
	r.GET(httpPath, func(c *gin.Context) {
		c.FileFromFS(fsPath, httpFS)
	})
}

func newWebIndexHandler(index []byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Write(index)
	}
}

func AddWebRoutes(r *gin.Engine, dist *embed.FS) {
	httpFS := http.FS(dist)

	// Route for /index.html
	if index, err := dist.ReadFile("web/dist/index.html"); err == nil {
		r.GET("/", newWebIndexHandler(index))
	} else {
		log.Fatal(err)
	}

	// Route for /*
	if files, err := fs.ReadDir(dist, "web/dist"); err == nil {
		for _, f := range files {
			name := f.Name()
			if !f.IsDir() && name != "index.html" {
				httpPath := "/" + name
				fsPath := path.Join("web/dist", name)
				newWebRouteHandler(r, httpFS, httpPath, fsPath)
			}
		}
	} else {
		log.Fatal(err)
	}

	// Route for /assets/*
	r.GET("/assets/*assets", func(c *gin.Context) {
		c.FileFromFS(path.Join("web/dist", c.Request.URL.Path), httpFS)
	})
}
