package auth

import (
	"context"
	"fmt"
	dto "lms_system/internal/domain/dto/auth"

	"github.com/sirupsen/logrus"
)

func (s *Service) UpdateProfile(ctx context.Context, userID string, request *dto.UpdateProfileRequest) (response *dto.UpdateProfileResponse, err error) {

	// INFO — общее действие
	s.logger.WithField("user_id", userID).Info("Updating user profile")

	// Get current user data from Keycloak
	user, err := s.keycloakClient.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user from keycloak")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update only the fields that were provided
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.Email != "" {
		user.Email = request.Email
	}

	// DEBUG — детали обновления
	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}).Debug("Profile update details")

	// Update user in Keycloak
	if err := s.keycloakClient.UpdateUser(ctx, userID, user); err != nil {
		s.logger.WithError(err).Error("Failed to update user in keycloak")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// INFO — успешный результат
	s.logger.WithField("user_id", userID).Info("User profile updated successfully")

	return &dto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	}, nil
}
