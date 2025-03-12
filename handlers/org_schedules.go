package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
	"github.com/rajnandan1/smaraka/validators"
)

// handle function fetch all org schedules
func (h *HandlersImplementation) GetOrgSchedules(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	schedules, err := h.db.GetAllSchedulesForORG(ctx, orgUser.OrganizationID)
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
	err := h.db.UpdateScheduleStatus(ctx, req.ScheduleID, orgUser.OrganizationID, req.Status)
	if err != nil {
		logger.LogError("Error updating org schedule", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	schedules, err := h.db.GetAllSchedulesForORG(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting org schedules", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, schedules)
}
func (h *HandlersImplementation) CreateOrgSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.CreateOrgScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	//validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	schedule := models.Schedule{
		ScheduleID:          h.db.NewID("schedule"),
		ScheduleName:        req.ScheduleName,
		ScheduleDescription: req.ScheduleDescription,
		ScheduleType:        req.ScheduleType,
		ScheduleStatus:      constants.ScheduleStatusActive,
		ScheduleURL:         req.ScheduleURL,
		OrganizationID:      orgUser.OrganizationID,
		ScheduleMeta:        "",
		IntervalDays:        req.IntervalDays,
	}

	err := h.db.InsertSchedule(ctx, &schedule)
	if err != nil {
		logger.LogError("Error updating org schedule", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	schedules, err := h.db.GetAllSchedulesForORG(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting org schedules", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, schedules)
}

func (h *HandlersImplementation) DeleteOrgSchedules(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.DeleteScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	//validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	err := h.db.DeleteScheduleByIDs(ctx, req.ScheduleIDs, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error deleting org schedule", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	schedules, err := h.db.GetAllSchedulesForORG(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting org schedules", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, schedules)

}
func (h *HandlersImplementation) RunOrgSchedules(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.DeleteScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	//validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}

	orgDataURLs, err := h.svc.PlaySchedule(ctx, req.ScheduleIDs, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error deleting org schedule", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	allURLs := make([]string, 0)
	for _, orgData := range *orgDataURLs {
		allURLs = append(allURLs, orgData.URLs...)
	}

	validURLs := make([]string, 0)
	for _, url := range allURLs {
		if validators.IsValidURL(url) {
			validURLs = append(validURLs, url)
		}
	}
	if len(validURLs) > 0 {
		validURLChunks := utils.ChunkStringArray(validURLs, h.config.MaxWorkers)

		for _, validURLsChunked := range validURLChunks {
			h.bg.SubmitURLs(ctx, validURLsChunked, orgUser.OrganizationID)
		}

	}

	go func() {
		for _, orgData := range *orgDataURLs {
			h.svc.BulkLightAndFullJob(orgData.URLs, orgData.OrganizationID)
		}
	}()

	return c.JSON(http.StatusOK, orgDataURLs)

}
