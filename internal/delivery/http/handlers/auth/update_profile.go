package auth

import (
	"encoding/json"
	"net/http"

	dto "lms_system/internal/domain/dto/auth"
	"lms_system/utils"
)

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Accepts FirstName, LastName, Email and updates user profile in Keycloak
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UpdateProfileRequest   true  "Profile data"
// @Success      200    {object}  dto.UpdateProfileResponse
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/profile [put]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	user := utils.GetUserFromContext(r.Context())
	if user == nil || user.UserID == "" {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var request dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateProfile(r.Context(), user.UserID, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
