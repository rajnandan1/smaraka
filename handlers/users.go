package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
)

func (h *HandlersImplementation) GetUserByID(c echo.Context) error {
	ctx := c.Request().Context()
	// Get the user from the context
	userToken := mddls.GetUserDBFromEchoContext(c)
	if userToken == nil {
		return c.JSON(http.StatusUnauthorized, models.Error{
			Message: constants.ERRORMSG_UNAUTHENTICATED,
			Code:    constants.ERRORCODE_UNAUTHENTICATED,
		})
	}

	userID := userToken.ID

	user, err := h.db.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORMSG_UNKNOWN_ERROR,
		})
	}
	return c.JSON(http.StatusOK, user)
}
