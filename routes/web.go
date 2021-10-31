package routes

import (
	"embed"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func AddWebRoutes(r *gin.Engine, dist *embed.FS) {
	httpfs := http.FS(dist)

	index, err := dist.ReadFile("web/dist/index.html")
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.Writer.Write(index)
	})

	r.GET("/assets/*assets", func(c *gin.Context) {
		c.FileFromFS(path.Join("web/dist", c.Request.URL.Path), httpfs)
	})

	r.GET("/apple-touch-icon.png", func(c *gin.Context) {
		c.FileFromFS("web/dist/apple-touch-icon.png", httpfs)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("web/dist/favicon.ico", httpfs)
	})

	r.GET("/manifest.webmanifest", func(c *gin.Context) {
		c.FileFromFS("web/dist/manifest.webmanifest", httpfs)
	})

	r.GET("/pwa-192x192.png", func(c *gin.Context) {
		c.FileFromFS("web/dist/pwa-192x192.png", httpfs)
	})

	r.GET("/pwa-512x512.png", func(c *gin.Context) {
		c.FileFromFS("web/dist/pwa-512x512.png", httpfs)
	})

	r.GET("/registerWS.js", func(c *gin.Context) {
		c.FileFromFS("web/dist/registerWS.js", httpfs)
	})

	r.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("web/dist/robots.txt", httpfs)
	})

	r.GET("/sw.js", func(c *gin.Context) {
		c.FileFromFS("web/dist/sw.js", httpfs)
	})

	r.GET("/workbox-afb9f189.js", func(c *gin.Context) {
		c.FileFromFS("web/dist/workbox-afb9f189.js", httpfs)
	})
}
