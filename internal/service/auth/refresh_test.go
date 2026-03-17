package auth

import (
	"context"
	"errors"
	"testing"

	dto "lms_system/internal/domain/dto/auth"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"
	"lms_system/internal/service/auth/mocks"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRefresh_Success(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	logger := logrus.New()

	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         logger,
	}

	mockKeycloak.On("RefreshToken", context.Background(), "test_token").Return(
		&model.TokenResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		},
		nil,
	)

	request := &dto.AuthRefreshRequest{RefreshToken: "test_token"}
	response, err := service.Refresh(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "new_access_token", response.AccessToken)
	assert.Equal(t, "new_refresh_token", response.RefreshToken)
	mockKeycloak.AssertExpectations(t)
}

func TestRefresh_KeycloakError(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	mockLogger := logrus.New()

	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         mockLogger,
	}

	mockKeycloak.On("RefreshToken", context.Background(), "invalid_token").Return(
		nil,
		errors.New("invalid token"),
	)

	request := &dto.AuthRefreshRequest{RefreshToken: "invalid_token"}
	response, err := service.Refresh(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockKeycloak.AssertExpectations(t)
}
