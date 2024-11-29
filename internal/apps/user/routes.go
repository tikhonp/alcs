package user

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/user/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/util/annalist"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func ConfigureUserGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, a annalist.Annalist) {
	sah := &handlers.UserHandler{Db: modelsFactory, Annalist: a}

    g.Use(auth.AuthRequiredMiddleware(modelsFactory.AuthUsers()))

	g.GET("", sah.MainPage)
}
