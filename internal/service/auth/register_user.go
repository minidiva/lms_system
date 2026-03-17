package auth

import (
	"context"
	"fmt"
	dto "lms_system/internal/domain/dto/auth"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"

	"github.com/sirupsen/logrus"
)

func (s *Service) RegisterUser(ctx context.Context, request *dto.UserRegistrationRequest) (response *dto.UserRegistrationResponse, err error) {

	// INFO — общее действие
	s.logger.WithField("username", request.Username).Info("Registering new user")

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

	// DEBUG — детали нового пользователя
	s.logger.WithFields(logrus.Fields{
		"username":   request.Username,
		"email":      request.Email,
		"first_name": request.FirstName,
		"last_name":  request.LastName,
		"roles":      request.Roles,
	}).Debug("New user details")

	// Create user in Keycloak
	userID, err := s.keycloakClient.CreateUser(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create user in keycloak")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// INFO — пользователь создан
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": request.Username,
	}).Info("User created successfully in Keycloak")

	// Assign roles to the user
	for _, role := range request.Roles {
		if err := s.keycloakClient.AssignRoleToUser(ctx, userID, role); err != nil {
			s.logger.WithError(err).Errorf("Failed to assign role %s to user %s", role, userID)
			// Continue even if role assignment fails
		}

		// DEBUG — детали назначения роли
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"role":    role,
		}).Debug("Role assigned to user")
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": request.Username,
		"email":    request.Email,
	}).Info("User registered successfully")

	return &dto.UserRegistrationResponse{
		UserID:   userID,
		Username: request.Username,
		Email:    request.Email,
	}, nil
}
