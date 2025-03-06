package postgres

import (
	"context"

	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) GetOrgSchedules(ctx context.Context, org_id string) ([]models.OrgScheduleResponse, error) {
	// Prepare SQL statement
	query := `
        SELECT 
            s.schedule_id, 
            s.schedule_name, 
            s.schedule_description, 
            s.schedule_url, 
            s.schedule_meta, 
            COALESCE(os.status, 'inactive') AS status,
            s.default_interval_days AS interval
        FROM 
            schedules s
        LEFT JOIN 
            org_schedules os ON s.schedule_id = os.schedule_id AND os.organization_id = $1
        ORDER BY 
            s.schedule_name
    `

	// Execute query
	rows, err := p.Pool.Query(ctx, query, org_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var schedules []models.OrgScheduleResponse
	for rows.Next() {
		var schedule models.OrgScheduleResponse
		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.Name,
			&schedule.Description,
			&schedule.URL,
			&schedule.Meta,
			&schedule.Status,
			&schedule.Interval,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// write a function to insert a new schedule, on conflict update status and interval
func (p *PostgresImplementation) InsertNewOrgSchedule(ctx context.Context, org_id string, schedule_id string, status string) error {
	query := `
				INSERT INTO org_schedules (organization_id, schedule_id, status)
				VALUES ($1, $2, $3)
				ON CONFLICT (organization_id, schedule_id) DO UPDATE
				SET status = $3
		`

	_, err := p.Pool.Exec(ctx, query, org_id, schedule_id, status)
	if err != nil {
		return err
	}

	return nil
}

// get daily schedule from org schedule that are active
func (p *PostgresImplementation) GetOrgSchedulesInterval(ctx context.Context, interval int) ([]models.OrgSchedule, error) {
	query := `
		SELECT 
			os.organization_id, 
			os.schedule_id
		FROM 
			org_schedules os
			INNER JOIN schedules s ON os.schedule_id = s.schedule_id
		WHERE 
			os.status = 'active' AND s.default_interval_days = $1
	`

	rows, err := p.Pool.Query(ctx, query, interval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.OrgSchedule
	for rows.Next() {
		var schedule models.OrgSchedule
		err := rows.Scan(
			&schedule.OrganizationID,
			&schedule.ScheduleID,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// get schedules for a particular schedule id from schedules table
func (p *PostgresImplementation) GetScheduleByID(ctx context.Context, schedule_id string) (*models.Schedule, error) {
	query := `
		SELECT 
			schedule_id, 
			schedule_name, 
			schedule_description, 
			schedule_url, 
			schedule_meta, 
			default_interval_days
		FROM 
			schedules
		WHERE 
			schedule_id = $1
	`

	row := p.Pool.QueryRow(ctx, query, schedule_id)
	var schedule models.Schedule
	err := row.Scan(
		&schedule.ScheduleID,
		&schedule.ScheduleName,
		&schedule.ScheduleDescription,
		&schedule.ScheduleURL,
		&schedule.ScheduleMeta,
		&schedule.DefaultIntervalDays,
	)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}
