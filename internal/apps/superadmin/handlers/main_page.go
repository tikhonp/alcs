package handlers

import (
	"net/http"
	"strconv"

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
    allOrganizations, err := sah.Db.AlcsOrganizations().GetAll() 
    if err != nil {
        return err
    }
    return util.TemplRender(c, views.Clients(user, allOrganizations))
}

func (sah *SuperAdminHandler) Client(c echo.Context) error {
    strId := c.Param("id")
    id, err := strconv.Atoi(strId)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound)
    }
    user, err := auth.GetUser(c, sah.Db.AuthUsers())
    if err != nil {
        return err
    }
    o, err := sah.Db.AlcsOrganizations().GetById(id)
    if err != nil {
        return err
    }
    return util.TemplRender(c, views.Client(user, o))
}

func (sah *SuperAdminHandler) CreateClientPage(c echo.Context) error {
    user, err := auth.GetUser(c, sah.Db.AuthUsers())
    if err != nil {
        return err
    }
    return util.TemplRender(c, views.CreateClient(user))
}
