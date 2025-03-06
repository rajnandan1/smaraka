package services

import (
	"context"
	"encoding/hex"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

func (s *ServicesImplementation) CreateNewSecret(ctx context.Context, userId, orgId, secretType, secretValue, secretName string) (*models.DbSecret, error) {

	dbSecretValue, err := s.cr.HashSecret(secretValue)
	if err != nil {
		return nil, err
	}

	secret := models.DbSecret{
		ID:             s.db.NewID("secret"),
		OrganizationID: orgId,
		SecretType:     secretType,
		SecretValue:    hex.EncodeToString(dbSecretValue),
		CurrentState:   constants.SecretStatusActive,
		CreatorID:      userId,
		SecretName:     secretName,
	}

	ierr := s.db.InsertNewSecret(ctx, secret)
	if ierr != nil {
		return nil, ierr
	}

	return &secret, nil
}

func (s *ServicesImplementation) GetSecretByValue(ctx context.Context, secretValue string) (*models.DbSecret, error) {
	dbSecretValue, err := s.cr.HashSecret(secretValue)
	if err != nil {
		return nil, err
	}
	return s.db.GetSecretByValue(ctx, hex.EncodeToString(dbSecretValue))
}
