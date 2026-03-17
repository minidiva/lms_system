package lesson

import (
	"encoding/json"
	"lms_system/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// GetLessonById godoc
// @Summary      Get lesson by ID
// @Description  Returns lesson by ID
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "Lesson ID"
// @Success      200  {object}  entity.Lesson
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/lessons/{id} [get]
func (h *Handler) GetLessonById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}

	userCtx := utils.GetUserFromContext(r.Context())
	if userCtx == nil || userCtx.UserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	lesson, err := h.service.GetLessonById(r.Context(), uint(id), userCtx.UserID)
	if err != nil {
		if err.Error() == "access denied" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lesson); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
