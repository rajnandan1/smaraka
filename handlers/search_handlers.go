package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
)

func (h *HandlersImplementation) SearchBookmarks(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.SearchBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	orgUser := mddls.GetOrgUserFromEchoContext(c)

	searchResult, err := h.db.SearchURLs(ctx, orgUser.OrganizationID, req.Needle, "")
	if err != nil {
		logger.LogError("Error searching bookmarks", err)
		return c.JSON(http.StatusOK, make([]models.URLResponses, 0))
	}
	if searchResult == nil {
		logger.LogInfo("No bookmarks found")
		return c.JSON(http.StatusOK, make([]models.URLResponses, 0))
	}

	//filter results with score more than 1
	bestScore := (searchResult)[0].Score
	passingScore := bestScore / float64(100/h.config.SearchAffinity)

	filteredResults := make([]*models.URLResponses, 0)
	for _, result := range searchResult {
		if result.Score >= passingScore {
			filteredResults = append(filteredResults, result)
		}
	}

	return c.JSON(http.StatusOK, filteredResults)

}

func (h *HandlersImplementation) JobQueueStatus(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	status, err := h.db.GetJobQueueStatusCount(ctx, orgUser.OrganizationID)
	if err != nil {
		logger.LogError("Error getting job queue status", err)
		return c.JSON(http.StatusOK, make(map[string]int))
	}

	pendingJobs, err := h.db.GetJobQueuePendingOlderThan1Hour(ctx, orgUser.OrganizationID)
	if err == nil && len(pendingJobs) > 0 {

		urls := make([]string, 0)
		for _, job := range pendingJobs {
			urls = append(urls, job.JobData)
			//update status to processing
			h.db.UpdateJobQueueStatus(ctx, orgUser.OrganizationID, job.JobData, constants.JobQueueStatusPending)
		}
		if len(urls) > 0 {
			validURLChunks := utils.ChunkStringArray(urls, h.config.MaxWorkers)
			for _, validURLsChunked := range validURLChunks {

				h.bg.SubmitURLs(ctx, validURLsChunked, orgUser.OrganizationID)
			}
		}
	}

	return c.JSON(http.StatusOK, status)
}

func (h *HandlersImplementation) GetAllBookmarks(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.GetBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Limit == 0 {
		req.Limit = 12
	}

	orgUser := mddls.GetOrgUserFromEchoContext(c)

	resp := models.URLListResponse{
		Data:   make([]*models.URLResponses, 0),
		NextID: "",
		IsLast: false,
	}

	//get all url orgs
	data, err := h.db.GetURLsForOrganization(ctx, orgUser.OrganizationID, req.NextID, req.Limit)
	if err != nil || data == nil {
		logger.LogError("Error getting url orgs", err)
		resp.IsLast = true
		return c.JSON(http.StatusOK, resp)
	}
	resp.NextID = data[len(data)-1].OrganizationRelationID
	resp.Data = data

	isLastData, err := h.db.GetURLsForOrganization(ctx, orgUser.OrganizationID, resp.NextID, 1)

	if err != nil || isLastData == nil || len(isLastData) == 0 {
		resp.IsLast = true

	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HandlersImplementation) GetBookmarkByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	orgUser := mddls.GetOrgUserFromEchoContext(c)

	bookmark, err := h.db.GetURLStoreByURLOrgIDOrgID(ctx, id, orgUser.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_BOOKMARK_NOT_FOUND,
			Code:    constants.ERRORCODE_BOOKMARK_NOT_FOUND,
		})
	}
	return c.JSON(http.StatusOK, bookmark)
}
func (h *HandlersImplementation) IndexBookmarkByID(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.PostIndexingRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	urlExists, err := h.db.GetURLOrganizationsByURLIDOrgID(ctx, req.ID, orgUser.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	bookmark, err := h.db.GetURLStoreByID(ctx, urlExists.URLID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_BOOKMARK_NOT_FOUND,
			Code:    constants.ERRORCODE_BOOKMARK_NOT_FOUND,
		})
	}
	urlStore, err := h.svc.GetContentEasy(bookmark.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	h.svc.DoContentCompleteByID(urlStore.URL)
	return c.JSON(http.StatusOK, bookmark)
}

func (h *HandlersImplementation) BulkDelete(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.BulkDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	deleteErr := h.db.DeleteURLsByIDs(ctx, req.IDs, orgUser.OrganizationID)
	if deleteErr != nil {
		logger.LogError("Error deleting bookmarks", deleteErr)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, nil)
}
func (h *HandlersImplementation) DeleteBookmarkByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	urlExists, err := h.db.GetURLOrganizationsByURLIDOrgID(ctx, id, orgUser.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	err = h.db.DeleteURLByID(ctx, urlExists.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, nil)
}
func (h *HandlersImplementation) ExportBookmarks(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	fileName := "bookmarks_okBookarms.html"
	createdAt := time.Now().Unix()
	lastModifiedAt := time.Now().Unix()
	htmlMiddle := ""
	//GetAllURLsForORg
	urlStores, urlOrgs, _ := h.db.GetAllURLsForORG(ctx, orgUser.OrganizationID)
	if urlOrgs != nil && len(*urlOrgs) > 0 {
		createdAt = (*urlOrgs)[0].CreatedAt.Unix()
		lastModifiedAt = (*urlOrgs)[len(*urlOrgs)-1].CreatedAt.Unix()
	}
	if urlStores != nil && len(*urlStores) > 0 {

		for i, urlStore := range *urlStores {
			correspondURLORG := (*urlOrgs)[i]
			createDate := correspondURLORG.CreatedAt.Unix()
			lmdate := correspondURLORG.UpdatedAt.Unix()
			htmlMiddle += `<DT><A HREF="` + urlStore.URL + `" ADD_DATE="` + strconv.FormatInt(createDate, 10) + `" LAST_MODIFIED="` + strconv.FormatInt(lmdate, 10) + `">` + urlStore.Title + `</A>`
		}
	}
	//create file in /tmp with html content like firefox export bookmarks
	htmlStart := `<!DOCTYPE>
			<HEAD>
			<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
			<TITLE>OkBookmarks</TITLE>
			</HEAD>
			<HTML>
			<BODY>
			<H1>OkBookmarks Export</H1>
			<DL>
				<DT>
					<H3 ADD_DATE="` + strconv.FormatInt(createdAt, 10) + `" LAST_MODIFIED="` + strconv.FormatInt(lastModifiedAt, 10) + `">Bookmarks</H3>
				</DT>
			</DL>`

	htmlEnd := `</BODY></HTML>`

	completeHTML := htmlStart + htmlMiddle + htmlEnd

	//create file in /tmp with html content like firefox export bookmarks
	folderLocation := "/tmp"
	fileLocation := folderLocation + "/" + fileName
	//create file
	file, err := os.Create(fileLocation)
	if err != nil {
		logger.LogError("Error creating file", err)
		return c.HTML(http.StatusInternalServerError, constants.ERRORMSG_UNKNOWN_ERROR)
	}
	defer file.Close()

	_, err = file.WriteString(completeHTML)
	if err != nil {
		logger.LogError("Error writing to file", err)
		return c.HTML(http.StatusInternalServerError, constants.ERRORMSG_UNKNOWN_ERROR)
	}

	//write to file
	//send file

	return c.Attachment(fileLocation, fileName)
}
