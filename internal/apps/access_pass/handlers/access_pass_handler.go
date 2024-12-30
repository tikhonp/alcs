package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/access_pass/views"
	"github.com/tikhonp/alcs/internal/db/models/alcs"
	"github.com/tikhonp/alcs/internal/db/utils"
	"github.com/tikhonp/alcs/internal/util"
)

type AccessPassRequest struct {
	VisitorName    string `form:"visitor_name" validate:"required"`
	VehicleNumber  string `form:"vehicle_number"`
	Purpose        string `form:"purpose" validate:"required"`
	ValidFrom      string `form:"valid_from" validate:"required"`
	ValidUntil     string `form:"valid_until" validate:"required"`
	ContactDetails string `form:"contact_details"`
	HostID         int    `form:"host_id" validate:"required"`
}

func (aph *AccessPassHandler) RequestPassForm(c echo.Context) error {
	organizationId, err := strconv.Atoi(c.QueryParam("organization_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid organization ID"})
	}
	hosts, err := aph.Db.AuthUsers().GetHostsForOrganization(organizationId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to fetch hosts"})
	}

	return util.TemplRender(c, views.AccessPassRequestForm(hosts))
}

func (aph *AccessPassHandler) SubmitPassRequest(c echo.Context) error {
	// Bind and validate the form input
	req := new(AccessPassRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	// Parse validity dates
	validFrom, err := time.Parse("2006-01-02", req.ValidFrom)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid valid_from date"})
	}
	validUntil, err := time.Parse("2006-01-02", req.ValidUntil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid valid_until date"})
	}

	// Save the access pass request
	pass := &alcs.AccessPass{
		VisitorName:    utils.NewNullString(req.VisitorName),
		VehicleNumber:  utils.NewNullString(req.VehicleNumber),
		Purpose:        utils.NewNullString(req.Purpose),
		ValidFrom:      validFrom,
		ValidUntil:     validUntil,
		ContactDetails: utils.NewNullString(req.ContactDetails),
		HostID:         req.HostID,
		Status:         "Pending",
	}
	err = aph.Db.AlcsAccessPasses().CreateAccessPass(pass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save access pass"})
	}

	// Notify approver via Telegram
	err = aph.Bot.NotifyApproverWithInlineButtons(pass, aph.Db.AlcsAccessPasses())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send notification"})
	}

	return c.Redirect(http.StatusSeeOther, "/user/access-passes")
}
