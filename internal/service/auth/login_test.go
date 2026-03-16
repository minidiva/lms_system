package auth

import (
	"context"
	"errors"
	"lms_system/internal/domain/dto"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"
	"lms_system/internal/service/auth/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin_Success(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	logger := logrus.New()
	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         logger,
	}
	mockKeycloak.On("GetToken", mock.Anything, "user@test.com", "password123").
		Return(&model.TokenResponse{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_123",
		}, nil)
	result, err := service.Login(context.Background(), &dto.AuthLoginRequest{
		Username: "user@test.com",
		Password: "password123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access_token_123", result.AccessToken)
	assert.Equal(t, "refresh_token_123", result.RefreshToken)
	mockKeycloak.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	logger := logrus.New()
	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         logger,
	}
	mockKeycloak.On("GetToken", mock.Anything, "user@test.com", "wrongpassword").
		Return(nil, errors.New("invalid credentials"))
	result, err := service.Login(context.Background(), &dto.AuthLoginRequest{
		Username: "user@test.com",
		Password: "wrongpassword",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "authentication failed")
	mockKeycloak.AssertExpectations(t)
}

func TestLogin_KeycloakUnavailable(t *testing.T) {
	mockKeycloak := new(mocks.MockKeycloakClient)
	logger := logrus.New()
	service := &Service{
		keycloakClient: mockKeycloak,
		logger:         logger,
	}
	mockKeycloak.On("GetToken", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("connection refused"))
	result, err := service.Login(context.Background(), &dto.AuthLoginRequest{
		Username: "user@test.com",
		Password: "password123",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	mockKeycloak.AssertExpectations(t)
}
