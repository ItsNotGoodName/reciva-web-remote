package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var dist embed.FS

func handleGetWeb(httpFS http.FileSystem, fsPath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.FileFromFS(fsPath, httpFS)
	}
}

func handleGetWebIndex(index []byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Write(index)
	}
}

func AddWebRoutes(r *gin.Engine) {
	httpFS := http.FS(dist)

	// GET /index.html
	if index, err := dist.ReadFile("dist/index.html"); err == nil {
		hw := handleGetWebIndex(index)
		r.GET("/", hw)
		r.GET("/index.html", hw)
	} else {
		log.Fatal("AddWebRoutes(ERROR):", err)
	}

	// GET /*
	if files, err := fs.ReadDir(dist, "dist"); err == nil {
		for _, f := range files {
			name := f.Name()
			if !f.IsDir() && name != "index.html" {
				httpPath := "/" + name
				fsPath := path.Join("dist", name)
				r.GET(httpPath, handleGetWeb(httpFS, fsPath))
			}
		}
	} else {
		log.Fatal("AddWebRoutes(ERROR):", err)
	}

	// GET /assets/*
	r.GET("/assets/*assets", func(c *gin.Context) {
		c.FileFromFS(path.Join("dist", c.Request.URL.Path), httpFS)
	})
}
