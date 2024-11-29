package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (sah *SuperAdminHandler) MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Yeah! Superadmin")
}
