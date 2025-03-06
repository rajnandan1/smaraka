package auth

type Auth interface {
}

type AuthImplementation struct {
}

func NewAuthImplementation(projectID string, secret string, sessionDuration int) (*AuthImplementation, error) {

	return &AuthImplementation{}, nil
}
