package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rajnandan1/smaraka/bg"
	"github.com/rajnandan1/smaraka/config"
	"github.com/rajnandan1/smaraka/crypt"
	"github.com/rajnandan1/smaraka/postgres"
	"github.com/rajnandan1/smaraka/services"
)

type Handlers interface {
	AddNewBookmark(c echo.Context) error
	SearchBookmarks(c echo.Context) error
	GetAllBookmarks(c echo.Context) error
	GetBookmarkByID(c echo.Context) error
	IndexBookmarkByID(c echo.Context) error
	GithubStarsImport(c echo.Context) error
	AddBulkNewBookmarks(c echo.Context) error
	GetBookmarkCount(c echo.Context) error
	FileUploadBrowsers(c echo.Context) error
	LandingPage(c echo.Context) error
	GetUserByID(c echo.Context) error
	PatchTextDataByID(c echo.Context) error
	JobQueueStatus(c echo.Context) error
	DeleteBookmarkByID(c echo.Context) error
	BulkDelete(c echo.Context) error
	TermsPage(c echo.Context) error
	RefundPage(c echo.Context) error
	ContactPage(c echo.Context) error
	PricingPage(c echo.Context) error
	ExportBookmarks(c echo.Context) error
	CreateNewExtensionKey(c echo.Context) error
	GetExtensionKeys(c echo.Context) error
	DeactivateKeys(c echo.Context) error
	BookmarkPresent(c echo.Context) error

	//sign up and login
	SignUp(c echo.Context) error
	CreateOrg(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error

	GetOrgSchedules(c echo.Context) error
	UpdateOrgSchedule(c echo.Context) error
	CreateOrgSchedule(c echo.Context) error
	DeleteOrgSchedules(c echo.Context) error
	RunOrgSchedules(c echo.Context) error
}
type HandlersImplementation struct {
	db     postgres.Postgres
	bg     bg.Background
	svc    services.Services
	crypto crypt.Crypt
	config config.Config
}

func ConfigureHandlers(db postgres.Postgres, bg bg.Background, svc services.Services, config config.Config, crypto crypt.Crypt) (Handlers, error) {
	return &HandlersImplementation{
		db:     db,
		bg:     bg,
		svc:    svc,
		config: config,
		crypto: crypto,
	}, nil
}

func NotFoundHandler(c echo.Context) error {
	return c.Render(http.StatusNotFound, "404.html", nil)
}
