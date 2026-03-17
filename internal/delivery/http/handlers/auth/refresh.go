package auth

import (
	"encoding/json"
	"net/http"

	dto "lms_system/internal/domain/dto/auth"
)

// Refresh godoc
// @Summary      User refresh
// @Description  Accepts RefreshKey, returns new AccessToken and RefreshToken
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.AuthRefreshRequest  true  "Refresh key"
// @Success      200    {object}  dto.AuthRefreshResponse
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request dto.AuthRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Refresh(r.Context(), &request)
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
