package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/user/views"
	db_auth "github.com/tikhonp/alcs/internal/db/models/auth"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (uh *UserHandler) MainPage(c echo.Context) error {
	user, err := auth.GetUser(c, uh.Db.AuthUsers())
	if err != nil {
		return err
	}
	hasSuperadminPermission, err := uh.Db.AuthUsers().IsUserHasPermissions(user.Id, db_auth.SuperAdmin)
    if err != nil {
        return err
    }
	return util.TemplRender(c, views.UserPage(user, hasSuperadminPermission))
}
