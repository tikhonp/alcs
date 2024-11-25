package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (ah *AuthHandler) OAuthCallback(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

	guser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx))
	if err != nil {
		return err
	}

	user, err := ah.Db.AuthUsers().FromOAuth(&guser)
	if err != nil {
		return err
	}

	auth.Login(c, user.Id)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (ah *AuthHandler) OAuthLogout(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

	gothic.Logout(c.Response(), c.Request().WithContext(ctx))
	c.Response().Header().Set("Location", "/")
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (ah *AuthHandler) OAuthProvider(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

	// try to get the user without re-authenticating
	if guser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx)); err == nil {
		user, err := ah.Db.AuthUsers().FromOAuth(&guser)
		if err != nil {
			return err
		}
		auth.Login(c, user.Id)
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	gothic.BeginAuthHandler(c.Response(), c.Request().WithContext(ctx))
	return nil
}
