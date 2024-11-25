package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/main_page/views"
	"github.com/tikhonp/alcs/internal/util"
)

func (mph *MainPageHandler) MainPage(c echo.Context) error {
    return util.TemplRender(c, views.MainPage())
}
