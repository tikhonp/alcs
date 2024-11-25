package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/auth/views"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (ah *AuthHandler) Login(c echo.Context) error {
	// Login page need to accept "next" get param with an url to redirect to
    nextPath := c.QueryParam("next")
    if nextPath == "" {
        nextPath = "/"
    }
    return util.TemplRender(c, views.LoginPage(nextPath)) 
}

func (ah *AuthHandler) Logout(c echo.Context) error {
    println("lol")
	auth.Logout(c)
	return c.Redirect(http.StatusMovedPermanently, "/")
}
