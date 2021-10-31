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

	r.GET("/assets/*assets", func(c *gin.Context) {
		c.FileFromFS(path.Join("web/dist", c.Request.URL.Path), httpfs)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("web/dist/favicon.ico", httpfs)
	})

	index, err := dist.ReadFile("web/dist/index.html")
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.Writer.Write(index)
	})
}
