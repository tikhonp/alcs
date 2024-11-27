package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/tikhonp/alcs/internal/apps/auth/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/util/annalist"
)

func ConfigureAuthGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, a annalist.Annalist) {
    ah := handlers.AuthHandler{Db: modelsFactory, Annalist: a}

	goth.UseProviders(
		google.New(cfg.Auth.GoogleKey, cfg.Auth.GoogleSecret, fmt.Sprintf("%s/auth/google/callback", cfg.BaseHost)),
	)

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return next(c)
		}
	})

	g.GET("/login", ah.LoginGet)
	g.POST("/login", ah.LoginPost)

	g.GET("/logout", ah.Logout)

	g.GET("/:provider/callback", ah.OAuthCallback)

	g.GET("/logout/:provider", ah.OAuthLogout)

	g.GET("/:provider", ah.OAuthProvider)

}
