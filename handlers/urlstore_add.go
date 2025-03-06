package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
	"github.com/rajnandan1/smaraka/validators"
)

func (h *HandlersImplementation) BookmarkPresent(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.CreateBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}
	if !validators.IsValidURL(req.URL) {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: constants.ERRORMSG_INVALID_URL,
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}

	orgUser := mddls.GetOrgUserFromEchoContext(c)
	oldURL, err := h.db.GetSingleURLForOrganizationURL(ctx, orgUser.OrganizationID, req.URL)
	if err != nil {
		logger.LogError("Error getting single url for org", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, oldURL)

}

func (h *HandlersImplementation) AddNewBookmark(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.CreateBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}
	if !validators.IsValidURL(req.URL) {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: constants.ERRORMSG_INVALID_URL,
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}

	orgUser := mddls.GetOrgUserFromEchoContext(c)

	_, err := h.db.GetURLStoreByURL(ctx, req.URL)

	if err != nil { //meaning new url, not present in db
		urlStore, err := h.svc.GetContentEasy(req.URL)
		if err != nil {
			logger.LogError("Error getting content", err)
			return c.JSON(http.StatusInternalServerError, models.Error{
				Message: constants.ERRORMSG_UNKNOWN_ERROR,
				Code:    constants.ERRORCODE_UNKNOWN_ERROR,
			})
		}
		urlStore.ID = h.db.NewID("url")
		if _, err := h.db.InsertNewURLStore(ctx, *urlStore); err != nil {
			logger.LogError("Error inserting url store", err)
			return c.JSON(http.StatusInternalServerError, models.Error{
				Message: constants.ERRORMSG_UNKNOWN_ERROR,
				Code:    constants.ERRORCODE_UNKNOWN_ERROR,
			})
		}

	}

	urlStore, err := h.db.GetURLStoreByURL(ctx, req.URL)
	if err != nil {
		logger.LogError("Error getting url store", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	//create new entry in url org
	urlOrg := models.URLOrganizations{
		ID:             h.db.NewID("url_org"),
		URLID:          urlStore.ID,
		OrganizationID: orgUser.OrganizationID,
		Status:         constants.URLStatusActive,
	}

	if _, err := h.db.InsertNewURLOrganization(ctx, urlOrg); err != nil {
		logger.LogError("Error inserting url org", err)
		if strings.Contains(err.Error(), "violates unique constraint") {
			oldURL, err := h.db.GetSingleURLForOrganizationURL(ctx, orgUser.OrganizationID, req.URL)
			if err != nil {
				logger.LogError("Error getting single url for org", err)
				return c.JSON(http.StatusInternalServerError, models.Error{
					Message: constants.ERRORMSG_UNKNOWN_ERROR,
					Code:    constants.ERRORCODE_UNKNOWN_ERROR,
				})
			}
			return c.JSON(http.StatusOK, oldURL)
		}
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	resp, err := h.db.GetSingleURLForOrganization(ctx, orgUser.OrganizationID, urlOrg.ID)
	if err != nil {
		logger.LogError("Error getting single url for org", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	go h.svc.DoContentCompleteByID(urlStore.URL)

	return c.JSON(http.StatusOK, resp)
}
func (h *HandlersImplementation) AddBulkNewBookmarks(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.CreateBulkBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}
	if req.Direction == "reverse" {
		//reverse the array req.URLs
		for i, j := 0, len(req.URLs)-1; i < j; i, j = i+1, j-1 {
			req.URLs[i], req.URLs[j] = req.URLs[j], req.URLs[i]
		}
	}
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	//create a new job
	validURLs := make([]string, 0)
	for _, url := range req.URLs {
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

	return c.JSON(http.StatusOK, map[string]string{
		"ok": "ok",
	})

}
