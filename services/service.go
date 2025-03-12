package services

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rajnandan1/smaraka/crypt"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/postgres"
)

type Services interface {
	ImportGithubStars(ghUrl string) (*[]models.GithubRepo, error)
	ParseUploadFile(fileObj models.FileUpload) ([]models.FileUploadResponse, error)
	GetContentEasy(url string) (*models.URLStore, error)
	DoContentCompleteByID(url_id string) (*models.URLStore, error)
	BulkLightAndFullJob(validURLs []string, orgId string) error
	CreateNewSecret(ctx context.Context, userId, orgId, secretType, secretValue, secretName string) (*models.DbSecret, error)
	GetSecretByValue(ctx context.Context, secretValue string) (*models.DbSecret, error)
	RunSchedule(ctx context.Context, interval int) (*[]models.PeriodicResponse, error)
	PlaySchedule(ctx context.Context, schedule_ids []string, org_id string) (*[]models.PeriodicResponse, error)
}
type ServicesImplementation struct {
	db     postgres.Postgres
	cr     crypt.Crypt
	policy *bluemonday.Policy
}

func ConfigureServices(db postgres.Postgres, c crypt.Crypt, p *bluemonday.Policy) (Services, error) {
	return &ServicesImplementation{
		db:     db,
		cr:     c,
		policy: p,
	}, nil
}
