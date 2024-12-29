package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/superadmin/views"
	"github.com/tikhonp/alcs/internal/db/models/alcs"
	"github.com/tikhonp/alcs/internal/db/utils"
	"github.com/tikhonp/alcs/internal/util"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (sah *SuperAdminHandler) CreateOrganizationPage(c echo.Context) error {
	user, err := auth.GetUser(c, sah.Db.AuthUsers())
	if err != nil {
		return err
	}
	return util.TemplRender(c, views.CreateOrganizationPage(user))
}

type organizationModel struct {
	Name  string `form:"name" validate:"required"`
	Notes string `form:"notes" validate:"required"`
}

func (sah *SuperAdminHandler) CreateOrganization(c echo.Context) error {
	// Bind and validate the form
	m := new(organizationModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}

	// Save the organization to the database
	org := alcs.Organization{
		Name:  m.Name,
		Notes: utils.NewNullString(m.Notes),
	}
	err := sah.Db.AlcsOrganizations().Create(&org)
	if err != nil {
		return err
	}

	// Redirect back to the list of organizations
	return c.Redirect(http.StatusSeeOther, "/superadmin/clients")
}
