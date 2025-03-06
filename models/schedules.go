package models

type PeriodicResponse struct {
	OrganizationID string    `json:"organization_id"`
	URLs           *[]string `json:"urls"`
}
