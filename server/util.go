package server

import "github.com/gin-gonic/gin"

func renderError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"err": err.Error()})
}
