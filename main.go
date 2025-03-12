package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rajnandan1/smaraka/bg"
	"github.com/rajnandan1/smaraka/config"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/crypt"
	"github.com/rajnandan1/smaraka/handlers"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/mddls"
	"github.com/rajnandan1/smaraka/migrations"
	"github.com/rajnandan1/smaraka/postgres"
	"github.com/rajnandan1/smaraka/services"
	"github.com/rajnandan1/smaraka/validators"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//load config
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatalf("error loading config: %v", configErr)
	}

	postgresConnectionString := config.GetPostgresURL()
	if config.Environment == constants.EnvDevelopment {
		log.Printf("postgresConnectionString: %s", postgresConnectionString)
	}

	crypto, err := crypt.ConfigureCrypt(ctx, config.VaultToken, config.SessionTimeout)
	if err != nil {
		log.Fatalf("error configuring crypt: %v", err)
	}
	htmlPolicy := bluemonday.UGCPolicy()
	// htmlPolicy.AllowElements("b", "strong", "p", "i", "em", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li", "blockquote", "code", "pre", "figure", "figcaption", "div", "span", "br", "hr", "td", "th", "tr", "table", "thead", "tbody", "tfoot", "caption")

	migrations.DoPostgresMigrationsUp(postgresConnectionString)

	logger.StartLogger(config.Environment)

	e := echo.New()
	e.HideBanner = true
	e.Validator = validators.NewValidator()
	// e.HidePort = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusNotFound {
				handlers.NotFoundHandler(c)
				return
			}
		}
		// Call the default HTTP error handler for other errors
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.Recover())

	e.Static("/app", "./build")

	postgresDb, err := postgres.ConfigurePostgres(ctx, postgresConnectionString)
	if err != nil {
		log.Fatalf("error configuring postgres: %v", err)
	}

	migrations.DoRiverMigrationUp(postgresDb)

	services, err := services.ConfigureServices(postgresDb, crypto, htmlPolicy)
	if err != nil {
		panic(err)
	}
	//configure bg
	bgjb, bgjbErr := bg.ConfigureBackground(ctx, postgresDb, services, config.MaxWorkers)
	if bgjbErr != nil {
		log.Fatalf("error configuring background: %v", bgjbErr)
	}

	handlers, err := handlers.ConfigureHandlers(postgresDb, bgjb, services, *config, crypto)
	if err != nil {
		panic(err)
	}

	//middlewares
	authMdl := mddls.AuthMiddleware(crypto, postgresDb)
	orgMdl := mddls.OrgIDMiddleware(postgresDb)

	e.POST("/api/ui/user/sign-up", handlers.SignUp)
	e.POST("/api/ui/user/login", handlers.Login)
	e.GET("/api/ui/logout", handlers.Logout)

	//requires auth
	e.POST("/api/ui/org/create", handlers.CreateOrg, authMdl)
	e.POST("/api/ui/url/new-bookmark", handlers.AddNewBookmark, authMdl, orgMdl)
	e.POST("/api/ui/url/search-bookmarks", handlers.SearchBookmarks, authMdl, orgMdl)
	e.GET("/api/ui/url/all-bookmarks", handlers.GetAllBookmarks, authMdl, orgMdl)

	e.GET("/api/ui/url/get-bookmark/:id", handlers.GetBookmarkByID, authMdl, orgMdl)
	e.DELETE("/api/ui/url/delete-bookmark/:id", handlers.DeleteBookmarkByID, authMdl, orgMdl)
	e.GET("/api/ui/url/get-bookmark-count", handlers.GetBookmarkCount, authMdl, orgMdl)
	e.PATCH("/api/ui/url/index-bookmark/:id", handlers.IndexBookmarkByID, authMdl, orgMdl)
	e.POST("/api/ui/url/import-github", handlers.GithubStarsImport, authMdl, orgMdl)
	e.POST("/api/ui/url/import-browsers", handlers.FileUploadBrowsers, authMdl, orgMdl)
	e.POST("/api/ui/url/import-bulk", handlers.AddBulkNewBookmarks, authMdl, orgMdl)
	e.POST("/api/ui/url/delete-bulk", handlers.BulkDelete, authMdl, orgMdl)
	e.PATCH("/api/ui/url/bookmark-update/:id", handlers.PatchTextDataByID, authMdl, orgMdl)
	e.GET("/api/ui/url/view-schedules", handlers.GetOrgSchedules, authMdl, orgMdl)
	e.PATCH("/api/ui/url/update-schedules", handlers.UpdateOrgSchedule, authMdl, orgMdl)
	e.PATCH("/api/ui/url/create-schedules", handlers.CreateOrgSchedule, authMdl, orgMdl)
	e.POST("/api/ui/url/delete-schedules", handlers.DeleteOrgSchedules, authMdl, orgMdl)
	e.POST("/api/ui/url/run-schedules", handlers.RunOrgSchedules, authMdl, orgMdl)

	e.GET("/api/ui/url/bookmarks-queue", handlers.JobQueueStatus, authMdl, orgMdl)
	e.GET("/api/ui/url/bookmarks-export", handlers.ExportBookmarks, authMdl, orgMdl)

	myFigure := figure.NewColorFigure("OkBookmarks", "doom", "yellow", true)
	myFigure.Print()

	go func() {
		if err := e.Start(":" + strconv.Itoa(config.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GracefulShutDownTimeout)*time.Second)
	defer cancel()
	postgresDb.Close()
	if config.Environment == constants.EnvProduction {
		if err := bgjb.Close(ctx); err != nil {
			log.Fatalf("error closing background: %v", err)
		}
	}

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	//e.Logger.Fatal(e.Start(":1323"))
}
