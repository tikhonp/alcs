package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/apps/auth/handlers"
)

func ConfigureAuthGroup(g *echo.Group) {
    ah := handlers.AuthHandler{}

    g.GET("/login", ah.Login)
}
