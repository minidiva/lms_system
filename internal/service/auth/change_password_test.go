package auth

import (
	"context"
	dto "lms_system/internal/domain/dto/auth"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"
	"lms_system/internal/service/auth/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangePassword_Success(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	logger := logrus.New()
	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         logger,
	}

	userID := "user-id-123"

	// Мок 1 — GetUserByID возвращает пользователя
	mockKeycloak.On("GetUserByID", mock.Anything, userID).
		Return(&model.UserRepresentation{
			Username: "user@test.com",
		}, nil)

	// Мок 2 — Login (проверка текущего пароля)
	mockKeycloak.On("GetToken", mock.Anything, "user@test.com", "oldpassword123").
		Return(&model.TokenResponse{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_123",
		}, nil)

	// Мок 3 — ChangePassword
	mockKeycloak.On("ChangePassword", mock.Anything, userID, "newpassword123").
		Return(nil)

	result, err := service.ChangePassword(context.Background(), userID, &dto.ChangePasswordRequest{
		CurrentPassword: "oldpassword123",
		NewPassword:     "newpassword123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Password changed successfully", result.Message)
	mockKeycloak.AssertExpectations(t)
}
