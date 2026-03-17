package auth

import (
	"context"
	"fmt"
	dto "lms_system/internal/domain/dto/auth"

	"github.com/sirupsen/logrus"
)

func (s *Service) ChangePassword(ctx context.Context, userID string, request *dto.ChangePasswordRequest) (response *dto.ChangePasswordResponse, err error) {

	// INFO — общее действие
	s.logger.WithField("user_id", userID).Info("Changing user password")

	// Get current user data from Keycloak
	user, err := s.keycloakClient.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user from keycloak")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// DEBUG — детали пользователя
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": user.Username,
	}).Debug("User details for password change")

	// Verify current password by attempting to login
	loginRequest := &dto.AuthLoginRequest{
		Username: user.Username,
		Password: request.CurrentPassword,
	}
	_, err = s.Login(ctx, loginRequest)
	if err != nil {
		s.logger.WithError(err).Error("Current password verification failed")
		return nil, fmt.Errorf("current password is incorrect")
	}

	// Change password in Keycloak
	if err := s.keycloakClient.ChangePassword(ctx, userID, request.NewPassword); err != nil {
		s.logger.WithError(err).Error("Failed to change password in keycloak")
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	// INFO — успешный результат
	s.logger.WithField("user_id", userID).Info("Password changed successfully")

	return &dto.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}
