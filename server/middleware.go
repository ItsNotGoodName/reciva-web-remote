package server

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
)

func ensureUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "uuid is required"})
			return
		}

		c.Set("uuid", uuid)

		c.Next()
	}
}

func ensureRadio(h *radio.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "uuid is required"})
			return
		}

		rd, err := h.GetRadio(uuid)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}

		c.Set("radio", rd)

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	// https://github.com/gin-contrib/cors/issues/29#issuecomment-397859488
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
