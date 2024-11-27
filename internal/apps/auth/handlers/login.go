package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/auth/views"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (ah *AuthHandler) LoginGet(c echo.Context) error {
	// Login page need to accept "next" get param with an url to redirect to
	nextPath := c.QueryParam("next")
	if nextPath == "" {
		nextPath = "/"
	}
	return util.TemplRender(c, views.LoginPage(nextPath, ""))
}

type loginModel struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (ah *AuthHandler) LoginPost(c echo.Context) error {
    m := new(loginModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return util.TemplRender(c, views.LoginPage("", err.Error()))
	}
	err := auth.LoginByEmailAndPassword(c, ah.Db.AuthUsers(), m.Email, m.Password)
	if err != nil {
		return util.TemplRender(c, views.LoginPage("", "Неверный логин или пароль"))
	}
	nextPath := c.QueryParam("next")
	if nextPath == "" {
		nextPath = "/"
	}
	return c.Redirect(http.StatusMovedPermanently, nextPath)
}

func (ah *AuthHandler) Logout(c echo.Context) error {
	auth.Logout(c)
	return c.Redirect(http.StatusMovedPermanently, "/")
}
