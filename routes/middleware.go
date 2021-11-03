package routes

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func ensureUUID(a *api.API) func(c *gin.Context) {
	return func(c *gin.Context) {
		uuid, ok := c.Params.Get("uuid")
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Set("uuid", uuid)
	}
}

func ensureID(c *gin.Context) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.Set("id", id)
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
