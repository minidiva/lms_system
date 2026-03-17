package domain

import (
	"context"
	dto "lms_system/internal/domain/dto/auth"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, request *dto.AuthLoginRequest) (response *dto.AuthLoginResponse, err error)
	Refresh(ctx context.Context, request *dto.AuthRefreshRequest) (response *dto.AuthRefreshResponse, err error)
	RegisterUser(ctx context.Context, request *dto.UserRegistrationRequest) (response *dto.UserRegistrationResponse, err error)
	UpdateProfile(ctx context.Context, userID string, request *dto.UpdateProfileRequest) (response *dto.UpdateProfileResponse, err error)
	ChangePassword(ctx context.Context, userID string, request *dto.ChangePasswordRequest) (response *dto.ChangePasswordResponse, err error)
}
