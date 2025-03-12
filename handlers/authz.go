package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/validators"
	"github.com/sethvargo/go-diceware/diceware"
)

func (h *HandlersImplementation) LandingPage(c echo.Context) error {
	// htmlContent, err := ioutil.ReadFile("./gohtml/index.html")
	// if err != nil {
	// 	return err
	// }

	// return c.HTML(http.StatusOK, string(htmlContent))
	return c.Render(http.StatusOK, "index.html", nil)
}
func (h *HandlersImplementation) TermsPage(c echo.Context) error {
	return c.Render(http.StatusOK, "terms.html", map[string]interface{}{
		"env": h.config.Environment,
	})
}
func (h *HandlersImplementation) RefundPage(c echo.Context) error {
	return c.Render(http.StatusOK, "refund.html", map[string]interface{}{
		"env": h.config.Environment,
	})
}
func (h *HandlersImplementation) ContactPage(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", map[string]interface{}{
		"env": h.config.Environment,
	})
}
func (h *HandlersImplementation) PricingPage(c echo.Context) error {
	return c.Render(http.StatusOK, "pricing.html", map[string]interface{}{
		"env": h.config.Environment,
	})
}

// signup handler
func (h *HandlersImplementation) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.SignUpRequest
	if err := c.Bind(&req); err != nil {
		logger.LogError("Error binding signup request", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		logger.LogError("Error validating signup request", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(err.Error()))
	}

	if !validators.IsValidPassword(req.Password) {
		logger.LogError("Invalid password strength", nil)

		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(constants.ERRORMSG_INVALID_PASSWORD))
	}

	//generate password hash
	hash, err := h.crypto.HashPassword(req.Password)
	if err != nil {
		logger.LogError("Error hashing password", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(err.Error()))
	}

	newUser, err := h.db.InsertNewUser(ctx, models.Users{
		ID:           h.db.NewID("user"),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		SeenAt:       time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		logger.LogError("Error inserting new user", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(err.Error()))
	}

	//generate token
	token, exp, err := h.crypto.GenerateToken(newUser.ID)
	if err != nil {
		logger.LogError("Error generating token", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(err.Error()))
	}

	// h.config.
	//create cookie

	lastUsedOrg, orgUserErr := h.LastUserOrg(ctx, newUser.ID)
	if orgUserErr != nil {
		logger.LogError("Error getting last user organization", orgUserErr)

		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/signup?error="+url.QueryEscape(orgUserErr.Error()))
	}

	cookieToken := mddls.CreateCookie(exp, token, h.config.Environment, mddls.TokenCookieName)
	cookieOrg := mddls.CreateCookie(exp, lastUsedOrg.OrganizationID, h.config.Environment, mddls.OrgCookieName)
	c.SetCookie(cookieToken)
	c.SetCookie(cookieOrg)

	return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/home")

}

func (h *HandlersImplementation) LastUserOrg(ctx context.Context, userID string) (*models.UserOrganizations, error) {
	lastUsedOrg, orgUserErr := h.db.GetLastUserOrganizationByUserID(ctx, userID)
	if orgUserErr != nil {
		logger.LogError("Error getting user organization", orgUserErr)
		orgName := "some org name"
		if list, err := diceware.Generate(3); err == nil {
			orgName = strings.Join(list, "-")
		}
		newOrgId := h.db.NewID("org")

		org := models.Organizations{
			ID:        newOrgId,
			Name:      orgName,
			CreatorID: userID,
		}

		orgCreateErr := h.db.InsertNewOrganization(ctx, org)
		if orgCreateErr != nil {
			logger.LogError("Error creating organization", orgCreateErr)
			return nil, orgCreateErr
		}

		newUserOrgID := h.db.NewID("user_org")
		newUserOrgMapping := models.UserOrganizations{
			ID:             newUserOrgID,
			UserID:         userID,
			OrganizationID: org.ID,
			Role:           constants.RoleAdmin,
		}

		userOrgCreateErr := h.db.InsertNewUserOrganization(ctx, newUserOrgMapping)
		if userOrgCreateErr != nil {
			logger.LogError("Error creating user organization", userOrgCreateErr)
			return nil, userOrgCreateErr
		}

		return &newUserOrgMapping, nil

	}

	return lastUsedOrg, nil
}

// create org
func (h *HandlersImplementation) CreateOrg(c echo.Context) error {
	userDB := mddls.GetUserDBFromEchoContext(c)

	ctx := c.Request().Context()
	var req models.CreateOrgRequest
	if err := c.Bind(&req); err != nil {
		logger.LogError("Error binding create org request", err)
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_URL,
		})
	}

	if err := c.Validate(req); err != nil {
		logger.LogError("Error validating create org request", err)
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INVALID_DETAIL,
		})
	}

	//create org
	orgID := h.db.NewID("org")
	err := h.db.InsertNewOrganization(ctx, models.Organizations{
		ID:        orgID,
		Name:      req.Name,
		CreatorID: userDB.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		logger.LogError("Error inserting new organization", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INTERNAL_DETAIL,
		})
	}

	//get org
	org, err := h.db.GetOrganizationByID(ctx, orgID)
	if err != nil {
		logger.LogError("Error getting organization by ID", err)
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
			Code:    constants.ERRORCODE_INTERNAL_DETAIL,
		})
	}

	newUserOrgID := h.db.NewID("user_org")
	newUserOrgMapping := models.UserOrganizations{
		ID:             newUserOrgID,
		UserID:         userDB.ID,
		OrganizationID: org.ID,
		Role:           constants.RoleAdmin,
	}
	userOrgCreateErr := h.db.InsertNewUserOrganization(ctx, newUserOrgMapping)
	if userOrgCreateErr != nil {
		logger.LogError("Error creating user organization", userOrgCreateErr)
		return c.Redirect(http.StatusFound, "/login")
	}

	return c.JSON(http.StatusOK, org)
}

// login handler
func (h *HandlersImplementation) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.LogError("Error binding login form data", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(err.Error()))
	}

	// Validate the bound data
	if err := c.Validate(&req); err != nil {
		logger.LogError("Error validating login form data", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(err.Error()))

	}

	hashedPassword, err := h.db.GetPasswordHashByEmail(ctx, req.Email)
	if err != nil {
		logger.LogError("Error getting user by email", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(constants.ERRORMSG_USER_NOT_FOUND))
	}

	//check password
	if err = h.crypto.ComparePassword([]byte(hashedPassword), req.Password); err != nil {
		logger.LogError("Invalid password", nil)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(constants.ERRORCODE_INVALID_PASSWORD))

	}

	//get user
	user, err := h.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.LogError("Error getting user by email", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(constants.ERRORMSG_USER_NOT_FOUND))

	}

	//generate token
	token, exp, err := h.crypto.GenerateToken(user.ID)
	if err != nil {
		logger.LogError("Error generating token", err)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(err.Error()))

	}

	//get last used org
	lastUsedOrg, orgUserErr := h.LastUserOrg(ctx, user.ID)
	if orgUserErr != nil {
		logger.LogError("Error getting last user organization", orgUserErr)
		return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login?error="+url.QueryEscape(err.Error()))

	}

	//create cookie
	tokenCookie := mddls.CreateCookie(exp, token, h.config.Environment, mddls.TokenCookieName)
	orgCookie := mddls.CreateCookie(exp, lastUsedOrg.OrganizationID, h.config.Environment, mddls.OrgCookieName)

	c.SetCookie(tokenCookie)
	c.SetCookie(orgCookie)

	return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/home")

}

// create a logout handler that clears the cookie
func (h *HandlersImplementation) Logout(c echo.Context) error {

	expiredTime := time.Now().Add(-1 * time.Hour)
	tokenCookie := mddls.CreateCookie(expiredTime, "", h.config.Environment, mddls.TokenCookieName)
	c.SetCookie(tokenCookie)
	return c.Redirect(http.StatusFound, h.config.FrontBasePath+"/login")
}
