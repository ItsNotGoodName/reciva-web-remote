package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func AddConfigRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/config.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, cfg)
	})

}
