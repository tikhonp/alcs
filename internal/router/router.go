package router

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/markbates/goth/gothic"
	"github.com/tikhonp/alcs/internal/apps/auth"
	mainpage "github.com/tikhonp/alcs/internal/apps/main_page"
	"github.com/tikhonp/alcs/internal/apps/superadmin"
	"github.com/tikhonp/alcs/internal/apps/user"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/annalist"
	authutil "github.com/tikhonp/alcs/internal/util/auth"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Debug = cfg.Server.Debug

	// TODO: Set proper logger
	e.Logger.SetLevel(log.DEBUG)

	e.Pre(middleware.RemoveTrailingSlash())

	if e.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: e.Logger.Output(),
		}))
	} else {
		e.Use(middleware.Logger())
	}
	// e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			// TODO: Set proper CORS configuration
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{echo.GET},
		},
	))

	store := sessions.NewCookieStore([]byte(cfg.Server.Secret))
	gothic.Store = store
	e.Use(session.Middleware(store))

	e.Use(authutil.AuthMiddleware())

	e.Validator = util.NewDefaultValidator()
    e.HTTPErrorHandler = util.HTTPErrorHandler

	return e
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, modelsFactory db.ModelsFactory, am annalist.AnnalistManager) {
	mainpage.ConfigureMainPageGroup(e.Group(""), cfg, modelsFactory)
	auth.ConfigureAuthGroup(e.Group("/auth"), cfg, modelsFactory, am.GetAnnalist("AUTH"))
	superadmin.ConfigureSuperAdminGroup(e.Group("/superadmin"), cfg, modelsFactory, am.GetAnnalist("SUPERADMIN"))
	user.ConfigureUserGroup(e.Group("/user"), cfg, modelsFactory, am.GetAnnalist("USER"))
}

func Start(e *echo.Echo, cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return e.Start(addr)
}
