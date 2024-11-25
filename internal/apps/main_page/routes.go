package mainpage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/main_page/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
)

func ConfigureMainPageGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory) {
    mpg := handlers.MainPageHandler{}

    g.GET("/", mpg.MainPage)
}
