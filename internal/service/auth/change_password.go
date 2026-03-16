package auth

import (
	"context"
	"fmt"
	"lms_system/internal/domain/dto"
)

func (s *Service) ChangePassword(ctx context.Context, userID string, request *dto.ChangePasswordRequest) (response *dto.ChangePasswordResponse, err error) {
	// Get current user data from Keycloak
	user, err := s.keycloakClient.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user from keycloak")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password by attempting to login
	loginRequest := &dto.AuthLoginRequest{
		Username: user.Username,
		Password: request.CurrentPassword,
	}

	_, err = s.Login(ctx, loginRequest)
	if err != nil {
		s.logger.WithError(err).Error("current password verification failed")
		return nil, fmt.Errorf("current password is incorrect")
	}

	// Change password in Keycloak
	if err := s.keycloakClient.ChangePassword(ctx, userID, request.NewPassword); err != nil {
		s.logger.WithError(err).Error("failed to change password in keycloak")
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	return &dto.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}
