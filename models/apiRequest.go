package models

type CreateBookmarkRequest struct {
	URL string `json:"url"`
}
type CreateBulkBookmarkRequest struct {
	URLs      []string `json:"urls"`
	Direction string   `json:"direction"`
}

type PatchBookmarkRequest struct {
	ID string `json:"id"`
}
type SearchBookmarkRequest struct {
	Needle string `json:"needle"`
}

type GetBookmarkRequest struct {
	Status    string `query:"status"`
	CreatedAt string `query:"created_at"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	NextID    string `query:"next_id"`
	FetchType string `query:"fetch_type"`
}

type PostIndexingRequest struct {
	ID string `param:"id"`
}
type PostGithubStarsImportRequest struct {
	Username string `json:"username"`
}

type PatchTextDataRequest struct {
	Title   string `json:"title"`
	Text    string `json:"full_text"`
	Excerpt string `json:"excerpt"`
	ID      string `param:"id"`
}

type URLListResponse struct {
	Data   []*URLResponses `json:"data"`
	NextID string          `json:"next_id"`
	IsLast bool            `json:"is_last"`
}

type URLResponses struct {
	URLID                  string  `json:"url_id"`
	Title                  string  `json:"title"`
	URL                    string  `json:"url"`
	Excerpt                string  `json:"excerpt"`
	ImageSmall             string  `json:"image_small"`
	ImageLarge             string  `json:"image_large"`
	AccentColor            string  `json:"accent_color"`
	OrganizationRelationID string  `json:"organization_relation_id"`
	OrganizationURLStatus  string  `json:"organization_url_status"`
	Checked                bool    `json:"checked"`
	Score                  float64 `json:"score"`
}

type BulkDeleteRequest struct {
	IDs []string `json:"organization_relation_ids"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type SecretCreateRequest struct {
	SecretName string `json:"secret_name"`
}
type SecretCreateResponse struct {
	SecretName  string `json:"secret_name"`
	SecretValue string `json:"secret_value"`
	SecretType  string `json:"secret_type"`
}

type SecretDeactivateRequest struct {
	SecretID string `json:"secret_id"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}
type SignUpResponse struct {
	Token string `json:"token"`
	OrgID string `json:"org_id"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

type CreateOrgRequest struct {
	Name string `json:"name" validate:"required"`
}
type CreateOrgResponse struct {
	Token string `json:"token"`
}

type UpdateOrgScheduleRequest struct {
	ScheduleID string `json:"schedule_id"`
	Status     string `json:"status"`
}

type OrgScheduleResponse struct {
	ScheduleID  string `json:"schedule_id"`
	Name        string `json:"schedule_name"`
	Description string `json:"schedule_description"`
	URL         string `json:"schedule_url"`
	Meta        string `json:"schedule_meta"`
	Status      string `json:"status"`
	Interval    int    `json:"interval"`
}
