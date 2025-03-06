package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
)

// handle function fetch all org schedules
func (h *HandlersImplementation) GetOrgSchedules(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	schedules, err := h.db.GetOrgSchedules(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting org schedules", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, schedules)
}

// handler function update org schedule
func (h *HandlersImplementation) UpdateOrgSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.UpdateOrgScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}
	err := h.db.InsertNewOrgSchedule(ctx, orgUser.OrganizationID, req.ScheduleID, req.Status)
	if err != nil {
		logger.LogError("Error updating org schedule", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	schedules, err := h.db.GetOrgSchedules(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting org schedules", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, schedules)
}
