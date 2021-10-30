package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

type ConfigJSON struct {
	PresetsEnabled bool `json:"presetsEnabled"`
}

func getConfigHandler(cfg *config.Config) func(c *gin.Context) {
	cfgJSON := ConfigJSON{}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, cfgJSON)
	}
}

func AddConfigRoutes(r *gin.RouterGroup, cfg *config.Config) {
	r.GET("/config", getConfigHandler(cfg))
}
