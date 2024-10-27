package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tikhonp/alcs/db/models/auth"
	"github.com/tikhonp/alcs/middleware"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := middleware.GetValue("userId", c)
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

func PermissionMiddleware(users auth.Users, permissionCodenames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, ok := c.Get("userId").(int)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			granted, err := users.IsUserHasPermissions(userId, permissionCodenames...)
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

