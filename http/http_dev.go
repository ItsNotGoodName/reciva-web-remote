//go:build dev

package http

import (
	_ "github.com/ItsNotGoodName/reciva-web-remote/docs/swagger" // docs is generated by Swag CLI, you have to import it.
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func swagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
