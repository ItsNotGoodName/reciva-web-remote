package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func ensureUUID(a *api.API) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, ok := c.Params.Get("UUID")
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
