package auth

import (
	"context"
	"fmt"
	"lms_system/internal/domain/dto"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"
)

func (s *Service) RegisterUser(ctx context.Context, request *dto.UserRegistrationRequest) (response *dto.UserRegistrationResponse, err error) {
	// Create user representation for Keycloak
	user := &model.UserRepresentation{
		Username:      request.Username,
		Email:         request.Email,
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		Enabled:       true,
		EmailVerified: false,
		Credentials: []model.CredentialRepresentation{
			{
				Type:      "password",
				Value:     request.Password,
				Temporary: false,
			},
		},
	}

	// Create user in Keycloak
	userID, err := s.keycloakClient.CreateUser(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("failed to create user in keycloak")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign roles to the user
	for _, role := range request.Roles {
		if err := s.keycloakClient.AssignRoleToUser(ctx, userID, role); err != nil {
			s.logger.WithError(err).Errorf("failed to assign role %s to user %s", role, userID)
			// Continue even if role assignment fails
		}
	}

	return &dto.UserRegistrationResponse{
		UserID:   userID,
		Username: request.Username,
		Email:    request.Email,
	}, nil
}
