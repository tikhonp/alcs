package superadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/superadmin/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	db_auth "github.com/tikhonp/alcs/internal/db/models/auth"
	"github.com/tikhonp/alcs/internal/util/annalist"
	"github.com/tikhonp/alcs/internal/util/auth"
)

// superadmin users have `superadmin` group permission
// that can be creates only by manage script
// i will write that script now

func ConfigureSuperAdminGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, a annalist.Annalist) {
	sah := &handlers.SuperAdminHandler{Db: modelsFactory, Annalist: a}

	g.Use(auth.PermissionMiddleware(modelsFactory.AuthUsers(), db_auth.SuperAdmin))

	g.GET("", sah.MainPage)
	g.GET("/clients", sah.Clients)
	g.GET("/clients/:id", sah.Client)
	g.GET("/clients/create", sah.CreateOrganizationPage)
	g.POST("/clients/create", sah.CreateOrganization)
}
