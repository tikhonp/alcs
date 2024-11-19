package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tikhonp/alcs/config"
	authutil "github.com/tikhonp/alcs/util/auth"
    "github.com/tikhonp/alcs/apps/auth"
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

	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			// TODO: Set proper CORS configuration
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{echo.GET},
		},
	))

    e.Use(authutil.AuthMiddleware())

	return e
}

func RegisterRoutes(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Купил мужик шляпу, а она ему как раз!")
	})
    
    auth.ConfigureAuthGroup(e.Group("/auth"))
}

func Start(e *echo.Echo, cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return e.Start(addr)
}

