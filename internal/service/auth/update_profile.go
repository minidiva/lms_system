package auth

import (
	"context"
	"fmt"
	"lms_system/internal/domain/dto"
)

func (s *Service) UpdateProfile(ctx context.Context, userID string, request *dto.UpdateProfileRequest) (response *dto.UpdateProfileResponse, err error) {
	// Get current user data from Keycloak
	user, err := s.keycloakClient.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user from keycloak")
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

	// Update user in Keycloak
	if err := s.keycloakClient.UpdateUser(ctx, userID, user); err != nil {
		s.logger.WithError(err).Error("failed to update user in keycloak")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	}, nil
}
