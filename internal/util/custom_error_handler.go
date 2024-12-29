package util

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	genericviews "github.com/tikhonp/alcs/internal/generic_views"
)

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	switch code {
	case http.StatusNotFound:
		err = TemplRenderWithCode(c, genericviews.Page404(), http.StatusNotFound)
		c.Logger().Error(fmt.Sprintf("HTTP %d: %v", code, err))
	default:
		c.Logger().Error(fmt.Sprintf("HTTP %d: %v", code, err))
	}

}
