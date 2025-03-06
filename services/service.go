package services

import (
	"context"

	"github.com/rajnandan1/smaraka/crypt"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/postgres"
)

type Services interface {
	ImportGithubStars(username string) (*[]models.GithubRepo, error)
	ParseUploadFile(fileObj models.FileUpload) ([]models.FileUploadResponse, error)
	GetContentEasy(url string) (*models.URLStore, error)
	DoContentCompleteByID(url_id string) (*models.URLStore, error)
	BulkLightAndFullJob(validURLs []string, orgId string) error
	CreateNewSecret(ctx context.Context, userId, orgId, secretType, secretValue, secretName string) (*models.DbSecret, error)
	GetSecretByValue(ctx context.Context, secretValue string) (*models.DbSecret, error)
	DailySchedules(ctx context.Context, interval int) (*[]models.PeriodicResponse, error)
}
type ServicesImplementation struct {
	db postgres.Postgres
	cr crypt.Crypt
}

func ConfigureServices(db postgres.Postgres, c crypt.Crypt) (Services, error) {
	return &ServicesImplementation{
		db: db,
		cr: c,
	}, nil
}
