package bg

import (
	"context"
	"fmt"

	"github.com/rajnandan1/smaraka/services"
	"github.com/riverqueue/river"
)

type URLStoreProcessArgs struct {
	URLs    []string `json:"urls"`
	OrgUser string   `json:"org_user"`
}

func (URLStoreProcessArgs) Kind() string { return "url_store_process" }

type URLStoreProcessWorker struct {
	river.WorkerDefaults[URLStoreProcessArgs]
	Service services.Services
}

func (w *URLStoreProcessWorker) Work(ctx context.Context, job *river.Job[URLStoreProcessArgs]) error {
	err := w.Service.BulkLightAndFullJob(job.Args.URLs, job.Args.OrgUser)
	fmt.Println("Job done")
	return err
}

type PeriodicJobArgs struct {
	Interval int `json:"interval"`
}

type PeriodicJobWorker struct {
	river.WorkerDefaults[PeriodicJobArgs]
	Service services.Services
}

func (PeriodicJobArgs) Kind() string { return "periodic" }

func (w *PeriodicJobWorker) Work(ctx context.Context, job *river.Job[PeriodicJobArgs]) error {

	fmt.Printf("Running %s job at %d\n", job.Kind, job.Args.Interval)

	orgDataURLs, err := w.Service.DailySchedules(ctx, job.Args.Interval)
	if err != nil {

		return err
	}

	//loop through the orgDataURLs and create a job for each
	for _, orgData := range *orgDataURLs {
		w.Service.BulkLightAndFullJob(*orgData.URLs, orgData.OrganizationID)
	}

	return nil
}
