//go:build dev

package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func swagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
