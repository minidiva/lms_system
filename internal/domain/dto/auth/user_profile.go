package dto

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name"  validate:"omitempty,min=2"`
	Email     string `json:"email"      validate:"omitempty,email"`
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
}
