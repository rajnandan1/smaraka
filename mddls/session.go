package mddls

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/auth"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/crypt"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/postgres"
	"github.com/rajnandan1/smaraka/services"
)

const (
	// CookieName is the name of the cookie used to store the session token
	TokenCookieName = "smaraka_session"
	OrgCookieName   = "smaraka_org"
)

// AuthMiddleware function to authenticate requests using your Auth interface
func AuthMiddleware(crypt crypt.Crypt, db postgres.Postgres) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()
			// Get the Authorization header
			cookie, err := c.Cookie(TokenCookieName)
			if err != nil {
				logger.LogError("Cookie not found in request", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}
			claims, errS := crypt.VerifyToken(cookie.Value)
			if errS != nil {
				logger.LogError("Invalid token in cookie", errS)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			userID := claims.UserID
			//check if claims has has not expired
			//@TODO

			//get user using email
			userDB, err := db.GetUserByID(ctx, userID)
			if err != nil {
				logger.LogError("Failed to get user by ID", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			//set user in context
			c.Set("userDB", userDB)

			return next(c)
		}
	}
}

func OrgIDMiddleware(db postgres.Postgres) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//get path param
			cookie, err := c.Cookie(OrgCookieName)
			if err != nil {
				logger.LogError("Cookie not found in request", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}
			orgID := cookie.Value
			if orgID == "" {
				logger.LogError("Missing org_id parameter", nil)
				return c.JSON(http.StatusBadRequest, models.Error{
					Message: constants.ERRORMSG_INVALID_ORG_ID,
					Code:    constants.ERRORCODE_INVALID_ORG_ID,
				})
			}

			userID := GetUserDBFromEchoContext(c).ID
			//get org user
			orgUser, err := db.GetUserOrganizationByOrgIDAndUserID(c.Request().Context(), orgID, userID)
			if err != nil {
				logger.LogError("User not authorized for this organization", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			c.Set("orgUser", orgUser)

			return next(c)
		}
	}
}
func AlreadyAuth(crypt crypt.Crypt) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			if cookie, err := c.Cookie(TokenCookieName); err == nil {
				_, errS := crypt.VerifyToken(cookie.Value)
				if errS == nil {
					//@TODO make /app in config
					return c.Redirect(http.StatusFound, "/app")
				}
				// Only log error if there was a token but it was invalid
				logger.LogError("Invalid token during AlreadyAuth check", errS)
			}

			// Proceed with the next handler if authentication is successful
			return next(c)
		}
	}
}

func GetUserDBFromEchoContext(c echo.Context) *models.Users {
	user := c.Get("userDB")
	if user == nil {
		return nil
	}
	if d, ok := user.(*models.Users); ok {
		return d
	}
	return nil
}
func GetOrgUserFromEchoContext(c echo.Context) *models.UserOrganizations {
	user := c.Get("orgUser")
	if user == nil {
		return nil
	}
	if d, ok := user.(*models.UserOrganizations); ok {
		return d
	}
	return nil
}
func CreateCookie(res time.Time, value string, env string, cookieName string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = value
	cookie.Path = "/"
	cookie.Expires = res
	//cookie should work in all localhost irrespective of port
	cookie.Domain = "localhost"
	if env == constants.EnvProduction {
		cookie.Domain = ".okbookmarks.com"
		cookie.HttpOnly = true
		cookie.Secure = true
	}

	return cookie
}

func CreateExpiredCookie(env string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = TokenCookieName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	//cookie should work in all localhost irrespective of port
	cookie.Domain = "localhost"
	if env == constants.EnvProduction {
		cookie.Domain = ".okbookmarks.com"
		cookie.HttpOnly = true
		cookie.Secure = true
	}

	return cookie
}

func AuthCorsMiddleware(authClient auth.Auth, db postgres.Postgres, svc services.Services) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				logger.LogError("Missing Authorization header", nil)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			//remove bearer
			token = token[7:]

			sec, err := svc.GetSecretByValue(ctx, token)
			if err != nil {
				logger.LogError("Invalid token value", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			userDB, err := db.GetUserByID(ctx, sec.CreatorID)
			if err != nil {
				logger.LogError("Failed to get user by ID from secret", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			orgDb, err := db.GetOrganizationByID(ctx, sec.OrganizationID)
			if err != nil {
				logger.LogError("Failed to get organization by ID from secret", err)
				return c.JSON(http.StatusUnauthorized, models.Error{
					Message: constants.ERRORMSG_UNAUTHENTICATED,
					Code:    constants.ERRORCODE_UNAUTHENTICATED,
				})
			}

			orgUser := &models.UserOrganizations{
				UserID:         userDB.ID,
				OrganizationID: orgDb.ID,
			}

			//set user in context
			c.Set("userDB", userDB)
			c.Set("orgUser", orgUser)

			return next(c)
		}
	}
}
