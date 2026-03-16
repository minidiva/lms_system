package mocks

import (
	"context"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"

	"github.com/stretchr/testify/mock"
)

type MockKeycloakClient struct {
	mock.Mock
}

func (m *MockKeycloakClient) GetToken(ctx context.Context, email, password string) (*model.TokenResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenResponse), args.Error(1)
}

func (m *MockKeycloakClient) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenResponse), args.Error(1)
}

func (m *MockKeycloakClient) CreateUser(ctx context.Context, user *model.UserRepresentation) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockKeycloakClient) AssignRoleToUser(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockKeycloakClient) GetUserByID(ctx context.Context, userID string) (*model.UserRepresentation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserRepresentation), args.Error(1)
}

func (m *MockKeycloakClient) UpdateUser(ctx context.Context, userID string, user *model.UserRepresentation) error {
	args := m.Called(ctx, userID, user)
	return args.Error(0)
}

func (m *MockKeycloakClient) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	args := m.Called(ctx, userID, newPassword)
	return args.Error(0)
}
