package auth

import (
	"context"
	"lms_system/internal/infrastructure/clients/http/keycloak/model"
)

type KeycloakClientInterface interface {
	GetToken(ctx context.Context, email, password string) (*model.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
	CreateUser(ctx context.Context, user *model.UserRepresentation) (string, error)
	AssignRoleToUser(ctx context.Context, userID string, role string) error
	GetUserByID(ctx context.Context, userID string) (*model.UserRepresentation, error)
	UpdateUser(ctx context.Context, userID string, user *model.UserRepresentation) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
}
