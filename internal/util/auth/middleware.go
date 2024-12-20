package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tikhonp/alcs/internal/db/models/auth"
	"github.com/tikhonp/alcs/internal/util"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := util.GetValue("userId", c)
			if err != nil {
				log.Infof("Error getting userId from session %s", err.Error())
				return next(c)
			}
			intUserId, ok := userId.(int)
			if ok {
				c.Set("userId", intUserId)
			}
			return next(c)
		}
	}
}

func AuthRequiredMiddleware(users auth.Users) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, ok := c.Get("userId").(int)
			if !ok {
				return c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			}
			return next(c)
		}
	}
}

func PermissionMiddleware(users auth.Users, permissions ...auth.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, ok := c.Get("userId").(int)
			if !ok {
				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}
			granted, err := users.IsUserHasPermissions(userId, permissions...)
			if err != nil {
				return err
			}
			if !granted {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			return next(c)
		}
	}
}
