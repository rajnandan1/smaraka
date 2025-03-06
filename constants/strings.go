package constants

const (
	// BookmarkStatusUnread is the status of an unread bookmark
	BookmarkStatusPending  = "PENDING"
	BookmarkStatusComplete = "COMPLETE"
	Browser                = "Browser"
	//Roles
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"

	//url status
	URLStatusActive  = "ACTIVE"
	URLStatusDeleted = "DELETED"

	//JobTypes
	JobTypeBookmark = "BOOKMARK"

	//JobStatus
	JobStatusPending  = "PENDING"
	JobStatusComplete = "COMPLETE"

	//JobQueue Status
	JobQueueStatusPending    = "PENDING"
	JobQueueStatusProcessing = "PROCESSING"
	JobQueueStatusComplete   = "COMPLETE"
	JobQueueStatusDeclined   = "DECLINED"
	JobQueueStatusFailed     = "FAILED"
	JobQueueStatusQueued     = "QUEUED"

	PrefixDatabaseUser    = "user"
	PrefixDatabaseOrg     = "org"
	PrefixDatabaseURL     = "url"
	PrefixDatabaseUserOrg = "user_org"
	PrefixDatabaseURLOrg  = "url_org"

	//HeadlessUserAgent
	HeadlessUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"

	//Envs
	EnvDevelopment = "development"
	EnvProduction  = "production"

	//SecretStatus
	SecretStatusActive   = "ACTIVE"
	SecretStatusInactive = "INACTIVE"

	//SecretTypes
	SecretTypeExtension = "EXTENSION"

	//SubscriptionType
	SubscriptionTypeInactive = "INACTIVE"
	SubscriptionTypeActive   = "ACTIVE"

	//OrderState
	OrderActive = "ACTIVE"
	OrderPaid   = "PAID"
)
