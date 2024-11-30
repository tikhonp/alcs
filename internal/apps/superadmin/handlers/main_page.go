package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/superadmin/views"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (sah *SuperAdminHandler) MainPage(c echo.Context) error {
    user, err := auth.GetUser(c, sah.Db.AuthUsers())
    if err != nil {
        return err
    }
    return util.TemplRender(c, views.SuperadminPage(user))
}

func (sah *SuperAdminHandler) Clients(c echo.Context) error {
    user, err := auth.GetUser(c, sah.Db.AuthUsers())
    if err != nil {
        return err
    }
    return util.TemplRender(c, views.Clients(user))
}
