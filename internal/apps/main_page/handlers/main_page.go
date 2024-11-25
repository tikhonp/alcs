package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/main_page/views"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (mph *MainPageHandler) MainPage(c echo.Context) error {
	user, err := auth.GetUser(c, mph.Db.AuthUsers())
	if errors.Is(err, auth.ErrNotAuthentificated) {
		return util.TemplRender(c, views.MainPage(false, user))
	} else if err != nil {
		return err
	}
	return util.TemplRender(c, views.MainPage(true, user))
}
