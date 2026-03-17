package dto

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}