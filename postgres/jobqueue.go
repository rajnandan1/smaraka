package postgres

import (
	"context"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) InsertJobQueues(ctx context.Context, org_id, job_id string, job_data []string) error {
	// Prepare SQL statement
	stmt := `INSERT INTO job_queue(id, org_id, job_id, job_data, status, created_at, updated_at) VALUES($1, $2, $3, $4, $5, now(), now())`
	for _, data := range job_data {
		_, err := p.Pool.Exec(ctx, stmt, p.NewID("jq"), org_id, job_id, data, constants.JobQueueStatusPending)
		if err != nil {
			return err
		}
	}
	return nil
}

// get jobqueue models.JobQueue by orgid and job data, should return only one
func (p *PostgresImplementation) GetJobQueueByOrgIDJobData(ctx context.Context, org_id, job_data string) (*models.JobQueue, error) {
	var jobQueue models.JobQueue
	query := `SELECT id, org_id, job_id, job_data, status, created_at, updated_at FROM job_queue WHERE org_id=$1 and job_data=$2`
	row := p.Pool.QueryRow(ctx, query, org_id, job_data)
	err := row.Scan(&jobQueue.ID, &jobQueue.OrgID, &jobQueue.JobID, &jobQueue.JobData, &jobQueue.Status, &jobQueue.CreatedAt, &jobQueue.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &jobQueue, nil
}

// insert one on unique key constraint violation update status to pending
func (p *PostgresImplementation) InsertJobQueue(ctx context.Context, org_id, job_data string) error {
	// Prepare SQL statement
	stmt := `INSERT INTO job_queue(id, org_id, job_id, job_data, status, created_at, updated_at) VALUES($1, $2, $3, $4, $5, now(), now()) ON CONFLICT (job_data, org_id) DO UPDATE SET status=$5, updated_at=now()`
	_, err := p.Pool.Exec(ctx, stmt, p.NewID("jq"), org_id, p.NewID("job"), job_data, constants.JobQueueStatusPending)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresImplementation) UpdateJobQueueStatus(ctx context.Context, org_id, job_data, status string) error {
	// Prepare SQL statement
	stmt := `UPDATE job_queue SET status=$1, updated_at=now() WHERE org_id=$2 and job_data=$3`
	_, err := p.Pool.Exec(ctx, stmt, status, org_id, job_data)
	if err != nil {
		return err
	}
	return nil
}

// GetJobQueueStatusCount(org_id string)
func (p *PostgresImplementation) GetJobQueueStatusCount(ctx context.Context, org_id string) (map[string]int, error) {
	var statusCount = make(map[string]int)
	query := `SELECT status, count(status) FROM job_queue WHERE org_id=$1 GROUP BY status`
	rows, err := p.Pool.Query(ctx, query, org_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		err = rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		statusCount[status] = count
	}
	return statusCount, nil
}

// given orgid get job_data with status pending or queued older than 1 hour
func (p *PostgresImplementation) GetJobQueuePendingOlderThan1Hour(ctx context.Context, org_id string) ([]models.JobQueue, error) {
	var jobQueues []models.JobQueue
	query := `SELECT id, org_id, job_id, job_data, status, created_at, updated_at FROM job_queue WHERE org_id=$1 and status in ($2, $3) and updated_at < now() - interval '1 hour'`
	//print query with values
	rows, err := p.Pool.Query(ctx, query, org_id, constants.JobQueueStatusPending, constants.JobQueueStatusQueued)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var jobQueue models.JobQueue
		err = rows.Scan(&jobQueue.ID, &jobQueue.OrgID, &jobQueue.JobID, &jobQueue.JobData, &jobQueue.Status, &jobQueue.CreatedAt, &jobQueue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobQueues = append(jobQueues, jobQueue)
	}
	return jobQueues, nil
}
