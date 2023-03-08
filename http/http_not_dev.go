//go:build !dev

package http

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/docs"
	"github.com/labstack/echo/v4"
)

func swagger(e *echo.Echo) {
	e.GET("/api/swagger.json", func(c echo.Context) error {
		return c.String(http.StatusOK, docs.SwaggerJSON)
	})
}
