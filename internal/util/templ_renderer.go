package util

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func TemplRender(c echo.Context, t templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}

func TemplRenderWithCode(c echo.Context, t templ.Component, statusCode int) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = statusCode
	return t.Render(c.Request().Context(), c.Response().Writer)
}
