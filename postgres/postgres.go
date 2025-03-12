package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rajnandan1/smaraka/models"
	"go.jetify.com/typeid"
)

// Postgres interface defines the methods for the PostgreSQL implementation
type Postgres interface {
	Close() // Add any other methods you want to define
	NewID(prefix string) string
	GetConnectionPool() *pgxpool.Pool

	//users
	InsertNewUser(ctx context.Context, user models.Users) (*models.Users, error)
	GetUserByEmail(ctx context.Context, email string) (*models.Users, error)
	GetUserByID(ctx context.Context, id string) (*models.Users, error)
	UpdateUserSeenAtByID(ctx context.Context, userID string) error
	GetPasswordHashByEmail(ctx context.Context, email string) (string, error)

	//organizations
	InsertNewOrganization(ctx context.Context, organization models.Organizations) error
	GetOrganizationsByCreatorID(creatorID string) ([]models.Organizations, error)
	GetOrganizationByID(ctx context.Context, id string) (*models.Organizations, error)

	//userorganizations
	InsertNewUserOrganization(ctx context.Context, userOrganization models.UserOrganizations) error
	GetLastUserOrganizationByUserID(ctx context.Context, userID string) (*models.UserOrganizations, error)
	UpdateUserOrganizationUpdatedAt(ctx context.Context, id string) error
	GetUserOrganizationByOrgIDAndUserID(ctx context.Context, orgID, userID string) (*models.UserOrganizations, error)

	//urlstore
	InsertNewURLStore(ctx context.Context, urlStore models.URLStore) (*models.URLStore, error)
	GetURLStoreByURL(ctx context.Context, url string) (*models.URLStore, error)
	GetURLStoreByID(ctx context.Context, id string) (*models.URLStore, error)
	UpdateURLStoreByURL(ctx context.Context, url string, urlData models.URLStore) (*models.URLStore, error)
	UpdateURLStoreByID(ctx context.Context, id string, urlData models.URLStore) (*models.URLStore, error)
	GetURLStoreByIDs(ctx context.Context, ids []string) ([]models.URLStore, error)
	GetURLsByIDs(ctx context.Context, ids []string) ([]models.URLStore, error)

	//jobqueue
	InsertJobQueues(ctx context.Context, org_id, job_id string, job_data []string) error
	UpdateJobQueueStatus(ctx context.Context, job_id, url_id, status string) error
	GetJobQueueStatusCount(ctx context.Context, org_id string) (map[string]int, error)
	InsertJobQueue(ctx context.Context, org_id, job_data string) error
	GetJobQueuePendingOlderThan1Hour(ctx context.Context, org_id string) ([]models.JobQueue, error)

	//urlorganizations
	InsertNewURLOrganization(ctx context.Context, urlOrganization models.URLOrganizations) (*models.URLOrganizations, error)
	GetURLsForOrganization(ctx context.Context, organizationID string, lastID string, pageSize int) ([]*models.URLResponses, error)
	GetURLOrganizationByID(ctx context.Context, id string) (*models.URLOrganizations, error)
	GetURLCountForOrganization(ctx context.Context, organizationID string) (int, error)
	GetURLOrganizationsByURLIDOrgID(ctx context.Context, urlID string, organizationID string) (*models.URLOrganizations, error)
	UpdateURLStatusByID(ctx context.Context, id string, status string) error
	DeleteURLByID(ctx context.Context, id string) error
	DeleteURLsByIDs(ctx context.Context, ids []string, orgID string) error
	GetNewURLsForOrganization(ctx context.Context, organizationID string, firstId string, pageSize int) ([]models.URLOrganizations, error)

	//urlstore and urlorganizations
	GetAllURLsForORG(ctx context.Context, orgID string) (*[]models.URLStore, *[]models.URLOrganizations, error)
	SearchURLs(ctx context.Context, orgID, query, domain string) ([]*models.URLResponses, error)
	GetURLStoreByURLOrgIDOrgID(ctx context.Context, urlOrgID, orgID string) (*models.URLStore, error)
	GetSingleURLForOrganization(ctx context.Context, organizationID string, urlOrgID string) (*models.URLResponses, error)
	GetSingleURLForOrganizationURL(ctx context.Context, organizationID string, url string) (*models.URLResponses, error)

	//secrets
	InsertNewSecret(ctx context.Context, secret models.DbSecret) error
	GetSecretByOrgAndValue(ctx context.Context, organizationID, secretType, secretValue string) (*models.DbSecret, error)
	GetSecretByID(ctx context.Context, id string) (*models.DbSecret, error)
	GetSecretsByOrganizationID(ctx context.Context, organizationID, secretType string) ([]*models.DbSecret, error)
	DeactivateSecret(ctx context.Context, id, organizationID string) error
	GetSecretByValue(ctx context.Context, secretValue string) (*models.DbSecret, error)

	//subscriptions
	InsertSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error)
	IsAllowedToUse(ctx context.Context, userID string) (bool, error)
	UpdateSubscriptionStatus(ctx context.Context, subID, orderID string) error

	//orders
	InsertOrder(ctx context.Context, order models.Orders) (*models.Orders, error)
	UpdateOrder(ctx context.Context, order_id, order_object string) (*models.Orders, error)

	GetActiveOrgSchedules(ctx context.Context, org_id string) (*[]models.Schedule, error)
	InsertSchedule(ctx context.Context, schedule *models.Schedule) error
	GetAllSchedulesForORG(ctx context.Context, orgID string) (*[]models.Schedule, error)
	GetScheduleByID(ctx context.Context, schedule_id string) (*models.Schedule, error)
	UpdateScheduleStatus(ctx context.Context, schedule_id, org_id, status string) error
	DeleteScheduleByIDs(ctx context.Context, schedule_ids []string, org_id string) error
	GetActiveSchedulesWithInterval(ctx context.Context, interval int) (*[]models.Schedule, error)
	GetSchedulesByIDsAndOrgIDs(ctx context.Context, schedule_ids []string, org_id string) (*[]models.Schedule, error)
}

// PostgresImplementation holds the connection pool
type PostgresImplementation struct {
	Pool *pgxpool.Pool
}

// GetConnectionPool returns the connection pool
func (p *PostgresImplementation) GetConnectionPool() *pgxpool.Pool {
	return p.Pool
}

// ConfigurePostgres initializes the PostgreSQL connection pool
func ConfigurePostgres(ctx context.Context, connString string) (Postgres, error) {

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &PostgresImplementation{Pool: pool}, nil // Use Pool instead of pool
}

// Close closes the database connection pool
func (p *PostgresImplementation) Close() {
	p.Pool.Close()
}

func (p *PostgresImplementation) NewID(prefix string) string {
	tid, err := typeid.WithPrefix(prefix)
	if err != nil {
		panic(err)
	}
	return tid.String()
}
