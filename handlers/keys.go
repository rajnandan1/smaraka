package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
	"github.com/segmentio/ksuid"
)

func (h *HandlersImplementation) CreateNewExtensionKey(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.SecretCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	secret := h.db.NewID("okbookmarks_secret") + ksuid.New().String()

	res, err := h.svc.CreateNewSecret(ctx, orgUser.UserID, orgUser.OrganizationID, constants.SecretTypeExtension, secret, req.SecretName)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	secretResponse := models.SecretCreateResponse{
		SecretName:  res.SecretName,
		SecretValue: secret,
		SecretType:  res.SecretType,
	}
	return c.JSON(http.StatusOK, secretResponse)
}

func (h *HandlersImplementation) GetExtensionKeys(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	secrets, err := h.db.GetSecretsByOrganizationID(ctx, orgUser.OrganizationID, constants.SecretTypeExtension)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	//make all secret value blank
	for i := range secrets {
		secrets[i].SecretValue = ""
	}
	if secrets == nil {
		secrets = make([]*models.DbSecret, 0)
	}
	return c.JSON(http.StatusOK, secrets)
}

func (h *HandlersImplementation) DeactivateKeys(c echo.Context) error {
	ctx := c.Request().Context()
	orgUser := mddls.GetOrgUserFromEchoContext(c)
	var req models.SecretDeactivateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	err := h.db.DeactivateSecret(ctx, req.SecretID, orgUser.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Error{
			Message: constants.ERRORMSG_UNKNOWN_ERROR,
			Code:    constants.ERRORCODE_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, nil)
}
