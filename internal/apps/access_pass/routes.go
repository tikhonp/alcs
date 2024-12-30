package accesspass

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/access_pass/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
)

func ConfigureAccessPassRoutes(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory) {
	aph := &handlers.AccessPassHandler{
		Db: modelsFactory,
	}

	g.GET("/request", aph.RequestPassForm)
	g.POST("/request", aph.SubmitPassRequest)
}
