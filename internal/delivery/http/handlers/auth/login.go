package auth

import (
	"encoding/json"
	"net/http"

	dto "lms_system/internal/domain/dto/auth"
)

// Login godoc
// @Summary      User login
// @Description  Accepts username and password, returns AccessToken and RefreshToken
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.AuthLoginRequest  true  "Login credentials"
// @Success      200    {object}  dto.AuthLoginResponse
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.AuthLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(r.Context(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
