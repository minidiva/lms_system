package auth

import (
	"encoding/json"
	"net/http"

	"lms_system/internal/domain/dto"
)

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
