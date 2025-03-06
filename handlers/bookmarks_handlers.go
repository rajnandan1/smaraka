package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
)

func (h *HandlersImplementation) PatchTextDataByID(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.PatchTextDataRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_URL,
		})
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
	if bookmark.Status == constants.BookmarkStatusPending {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: constants.ERRORMSG_BOOKMARK_PENDING,
			Code:    constants.ERRORCODE_BOOKMARK_PENDING,
		})
	}
	newBookmark := models.URLStore{
		ID:          bookmark.ID,
		URL:         bookmark.URL,
		Title:       req.Title,
		ImageSmall:  bookmark.ImageSmall,
		ImageLarge:  bookmark.ImageLarge,
		Excerpt:     req.Excerpt,
		AccentColor: bookmark.AccentColor,
		FullText:    req.Text,
		Status:      bookmark.Status,
		CreatedAt:   bookmark.CreatedAt,
	}
	if _, err := h.db.UpdateURLStoreByID(ctx, req.ID, newBookmark); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, newBookmark)
}

func (h *HandlersImplementation) GithubStarsImport(c echo.Context) error {
	var req models.PostGithubStarsImportRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	stars, err := h.svc.ImportGithubStars(req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	return c.JSON(http.StatusOK, stars)
}

func (h *HandlersImplementation) GetBookmarkCount(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	count, err := h.db.GetURLCountForOrganization(ctx, orgUser.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, map[string]int{"count": count})
}

// write a handler function for file upload
func (h *HandlersImplementation) FileUploadBrowsers(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_File,
		})
	}

	//get file extension
	fileExtension := filepath.Ext(file.Filename)
	uploadFile := models.FileUpload{
		FileName:      file.Filename,
		FileSize:      file.Size,
		FileType:      file.Header.Get("Content-Type"),
		FileExtension: fileExtension,
		ImportType:    constants.Browser,
		FileContent:   "",
	}

	//for firefox file type is html
	if fileExtension != ".html" || uploadFile.FileType != "text/html" {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: constants.ERRORMSG_INVALID_File,
			Code:    constants.ERRORCODE_INVALID_File,
		})
	}

	//get file content as string
	src, err := file.Open()
	if err != nil {

		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	defer src.Close()

	// Read the file
	// Create a buffer to write the file content
	fileContent := make([]byte, file.Size)
	// Read the file content to the buffer
	_, err = src.Read(fileContent)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	uploadFile.FileContent = string(fileContent)

	urls, err := h.svc.ParseUploadFile(uploadFile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}

	return c.JSON(http.StatusOK, urls)
}
