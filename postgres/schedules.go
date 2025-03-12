package postgres

import (
	"context"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

// insert a new schedule
func (p *PostgresImplementation) InsertSchedule(ctx context.Context, schedule *models.Schedule) error {
	query := `
		INSERT INTO schedules (schedule_id, schedule_name, schedule_description,schedule_type,schedule_status, schedule_url, schedule_meta, interval_days,organization_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := p.Pool.Exec(ctx, query, schedule.ScheduleID, schedule.ScheduleName, schedule.ScheduleDescription, schedule.ScheduleType, schedule.ScheduleStatus, schedule.ScheduleURL, schedule.ScheduleMeta, schedule.IntervalDays, schedule.OrganizationID)
	if err != nil {
		return err
	}

	return nil
}

// Get all schedules for an organization
func (p *PostgresImplementation) GetAllSchedulesForORG(ctx context.Context, orgID string) (*[]models.Schedule, error) {
	query := `
		SELECT 
			schedule_id,
			schedule_name,
			schedule_description,
			schedule_type,
			schedule_status,
			schedule_url,
			organization_id,
			schedule_meta,
			interval_days
		FROM 
			schedules
		WHERE 
			organization_id = $1
	`

	rows, err := p.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.ScheduleName,
			&schedule.ScheduleDescription,
			&schedule.ScheduleType,
			&schedule.ScheduleStatus,
			&schedule.ScheduleURL,
			&schedule.OrganizationID,
			&schedule.ScheduleMeta,
			&schedule.IntervalDays,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &schedules, nil
}

// Get Active Schedules for org
func (p *PostgresImplementation) GetActiveOrgSchedules(ctx context.Context, org_id string) (*[]models.Schedule, error) {
	// Prepare
	query := `
		SELECT
			schedule_id,
			schedule_name,
			schedule_description,
			schedule_type,
			schedule_status,
			schedule_url,
			organization_id,
			schedule_meta,
			interval_days
		FROM
			schedules
		WHERE
			organization_id = $1
			AND schedule_status = $2
	`

	// Execute
	rows, err := p.Pool.Query(ctx, query, org_id, constants.ScheduleStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse
	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.ScheduleName,
			&schedule.ScheduleDescription,
			&schedule.ScheduleType,
			&schedule.ScheduleStatus,
			&schedule.ScheduleURL,
			&schedule.OrganizationID,
			&schedule.ScheduleMeta,
			&schedule.IntervalDays,
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

	return &schedules, nil
}

// get schedules for a particular schedule id from schedules table
func (p *PostgresImplementation) GetScheduleByID(ctx context.Context, schedule_id string) (*models.Schedule, error) {
	query := `
		SELECT 
			schedule_id,
			schedule_name,
			schedule_description,
			schedule_type,
			schedule_status,
			schedule_url,
			organization_id,
			schedule_meta,
			interval_days
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
		&schedule.ScheduleType,
		&schedule.ScheduleStatus,
		&schedule.ScheduleURL,
		&schedule.OrganizationID,
		&schedule.ScheduleMeta,
		&schedule.IntervalDays,
	)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

// Update status of a schedule given schedule id and org id
func (p *PostgresImplementation) UpdateScheduleStatus(ctx context.Context, schedule_id, org_id, status string) error {
	query := `
		UPDATE schedules
		SET schedule_status = $1
		WHERE schedule_id = $2 AND organization_id = $3
	`

	_, err := p.Pool.Exec(ctx, query, status, schedule_id, org_id)
	if err != nil {
		return err
	}

	return nil
}

// Delete a schedule given schedule id and org id
func (p *PostgresImplementation) DeleteScheduleByIDs(ctx context.Context, schedule_ids []string, org_id string) error {
	query := `
		DELETE FROM schedules
		WHERE organization_id = $1 AND schedule_id = ANY($2)
	`

	_, err := p.Pool.Exec(ctx, query, org_id, schedule_ids)
	if err != nil {
		return err
	}

	return nil
}

// Get all schedules that are active and given an interval
func (p *PostgresImplementation) GetActiveSchedulesWithInterval(ctx context.Context, interval int) (*[]models.Schedule, error) {

	query := `
		SELECT
			schedule_id,
			schedule_name,
			schedule_description,
			schedule_type,
			schedule_status,
			schedule_url,
			organization_id,
			schedule_meta,
			interval_days
		FROM
			schedules
		WHERE
			schedule_status = $1
			AND interval_days = $2
	`

	rows, err := p.Pool.Query(ctx, query, constants.ScheduleStatusActive, interval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.ScheduleName,
			&schedule.ScheduleDescription,
			&schedule.ScheduleType,
			&schedule.ScheduleStatus,
			&schedule.ScheduleURL,
			&schedule.OrganizationID,
			&schedule.ScheduleMeta,
			&schedule.IntervalDays,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &schedules, nil
}

// get schedules for given schedule ids []string and org id
func (p *PostgresImplementation) GetSchedulesByIDsAndOrgIDs(ctx context.Context, schedule_ids []string, org_id string) (*[]models.Schedule, error) {
	query := `
		SELECT 
			schedule_id,
			schedule_name,
			schedule_description,
			schedule_type,
			schedule_status,
			schedule_url,
			organization_id,
			schedule_meta,
			interval_days
		FROM 
			schedules
		WHERE 
			organization_id = $1 AND schedule_id = ANY($2)
	`

	rows, err := p.Pool.Query(ctx, query, org_id, schedule_ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.ScheduleName,
			&schedule.ScheduleDescription,
			&schedule.ScheduleType,
			&schedule.ScheduleStatus,
			&schedule.ScheduleURL,
			&schedule.OrganizationID,
			&schedule.ScheduleMeta,
			&schedule.IntervalDays,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &schedules, nil
}
