package auth

import (
	"context"
	"fmt"
	dto "lms_system/internal/domain/dto/auth"
)

func (s *Service) Login(ctx context.Context, request *dto.AuthLoginRequest) (response *dto.AuthLoginResponse, err error) {

	// INFO — общее действие
	s.logger.Info("User login attempt")

	// DEBUG — детали входящего запроса
	s.logger.WithField("username", request.Username).Debug("Login request details")

	// Get token from Keycloak
	tokenResponse, err := s.keycloakClient.GetToken(ctx, request.Username, request.Password)
	if err != nil {
		// ERROR — уже есть ✅
		s.logger.WithError(err).Error("Failed to get token from keycloak")
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// INFO — успешный результат
	s.logger.WithField("username", request.Username).Info("User logged in successfully")

	return &dto.AuthLoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}
