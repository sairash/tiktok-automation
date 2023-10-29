package router

import (
	"beepbop/helper"
	api "beepbop/router/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Setuprouter(e *echo.Echo) {
	api.Api(e)

	e.GET("/", home)
}

func home(c echo.Context) error {
	helper.GetIPInfo(c)

	return c.String(http.StatusOK, "Hello, Home!")
}
