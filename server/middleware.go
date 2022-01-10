package server

import (
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
)

func ensureRadio(hub *radio.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			renderError(c, http.StatusBadRequest, fmt.Errorf("uuid is required"))
			c.Abort()
			return
		}

		rd, err := hub.GetRadio(uuid)
		if err != nil {
			renderError(c, http.StatusNotFound, err)
			c.Abort()
			return
		}

		c.Set("radio", rd)
	}
}
