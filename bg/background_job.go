package bg

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rajnandan1/smaraka/postgres"
	"github.com/rajnandan1/smaraka/services"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivershared/util/slogutil"
	"github.com/riverqueue/river/rivertype"
)

type Background interface {
	SubmitURLs(ctx context.Context, urls []string, orgId string) (*rivertype.JobInsertResult, error)
	Close(ctx context.Context) error
}

type BackgroundImplementation struct {
	riverClient  *river.Client[pgx.Tx]
	URLQueueName string
}

func ConfigureBackground(ctx context.Context, pg postgres.Postgres, svc services.Services, maxWorkers int) (Background, error) {

	queueName := "url_fetch"

	dbPool := pg.GetConnectionPool()
	workers := river.NewWorkers()

	// Register both workers
	river.AddWorker(workers, &URLStoreProcessWorker{
		Service: svc,
	})
	river.AddWorker(workers, &PeriodicJobWorker{
		Service: svc,
	}) // Add this line

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			queueName:          {MaxWorkers: maxWorkers},
			river.QueueDefault: {MaxWorkers: maxWorkers},
		},
		Logger:  slog.New(&slogutil.SlogMessageOnlyHandler{Level: slog.LevelWarn}),
		Workers: workers,
		PeriodicJobs: []*river.PeriodicJob{
			river.NewPeriodicJob(
				river.PeriodicInterval(24*time.Hour),
				func() (river.JobArgs, *river.InsertOpts) {
					return PeriodicJobArgs{
						Interval: 1,
					}, nil
				},
				&river.PeriodicJobOpts{RunOnStart: true},
			),
			river.NewPeriodicJob(
				river.PeriodicInterval(7*24*time.Hour),
				func() (river.JobArgs, *river.InsertOpts) {
					return PeriodicJobArgs{
						Interval: 7,
					}, nil
				},
				&river.PeriodicJobOpts{RunOnStart: true},
			), river.NewPeriodicJob(
				river.PeriodicInterval(30*24*time.Hour),
				func() (river.JobArgs, *river.InsertOpts) {
					return PeriodicJobArgs{
						Interval: 30,
					}, nil
				},
				&river.PeriodicJobOpts{RunOnStart: true},
			),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create River client: %v", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start River client: %v", err)
	}

	return &BackgroundImplementation{
		riverClient:  riverClient,
		URLQueueName: queueName,
	}, nil
}

func (b *BackgroundImplementation) Close(ctx context.Context) error {
	return b.riverClient.StopAndCancel(ctx)
}

func (b *BackgroundImplementation) SubmitURLs(ctx context.Context, urls []string, orgId string) (*rivertype.JobInsertResult, error) {
	res, err := b.riverClient.Insert(ctx, URLStoreProcessArgs{
		URLs:    urls,
		OrgUser: orgId,
	}, &river.InsertOpts{
		Queue:       b.URLQueueName,
		MaxAttempts: 3,
	})

	return res, err
}
