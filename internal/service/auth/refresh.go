package auth

import (
	"context"
	"fmt"
	dto "lms_system/internal/domain/dto/auth"
)

func (s *Service) Refresh(ctx context.Context, request *dto.AuthRefreshRequest) (response *dto.AuthRefreshResponse, err error) {
	// INFO — общее действие
	s.logger.Info("User refresh attempt")

	// DEBUG — детали входящего запроса
	s.logger.WithField("refresh_token", request.RefreshToken).Debug("Refresh request details")

	// Refresh token with Keycloak
	tokenResponse, err := s.keycloakClient.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		s.logger.WithError(err).Error("failed to refresh token from keycloak")
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	return &dto.AuthRefreshResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}
