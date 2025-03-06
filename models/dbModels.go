package models

import "time"

type URLStore struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Domain      string    `json:"domain"`
	Excerpt     string    `json:"excerpt"`
	ImageSmall  string    `json:"image_small"`
	ImageLarge  string    `json:"image_large"`
	Status      string    `json:"status"`
	FullText    string    `json:"full_text"`
	AccentColor string    `json:"accent_color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Score       float64   `json:"score"`
}
type Users struct {
	ID           string    `json:"id"`
	PasswordHash string    `json:"password_hash"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	SeenAt       time.Time `json:"seen_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Organizations struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatorID string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserOrganizations struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	OrganizationID string    `json:"organization_id"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type URLOrganizations struct {
	ID             string    `json:"id"`
	URLID          string    `json:"url_id"`
	OrganizationID string    `json:"organization_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobQueue struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	JobID     string    `json:"job_id"`
	JobData   string    `json:"job_data"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DbSecret struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	SecretType     string    `json:"secret_type"`
	SecretValue    string    `json:"secret_value"`
	SecretName     string    `json:"secret_name"`
	CurrentState   string    `json:"current_state"`
	CreatorID      string    `json:"creator_id"`
	CreatedAt      time.Time `json:"created_at"`
	LastUsedAt     time.Time `json:"last_used_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Subscription struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	SuccessOrderID string    `json:"success_order_id"`
	CurrentState   string    `json:"current_state"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type Orders struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	CustomerName   string    `json:"customer_name"`
	CustomerEmail  string    `json:"customer_email"`
	SubscriptionID string    `json:"subscription_id"`
	VariantID      int       `json:"variant_id"`
	StoreID        int       `json:"store_id"`
	CurrentState   string    `json:"current_state"`
	OrderObject    string    `json:"order_object"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type Schedule struct {
	ScheduleID          string    `json:"schedule_id"`
	ScheduleName        string    `json:"schedule_name"`
	ScheduleDescription string    `json:"schedule_description"`
	ScheduleURL         string    `json:"schedule_url"`
	ScheduleMeta        string    `json:"schedule_meta"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DefaultIntervalDays int       `json:"default_interval_days"`
}

type OrgSchedule struct {
	OrganizationID string    `json:"organization_id"`
	ScheduleID     string    `json:"schedule_id"`
	Status         string    `json:"status"`
	IntervalDays   int       `json:"interval_days"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
