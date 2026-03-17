package lesson

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteLessonById godoc
// @Summary      Delete lesson
// @Description  Deletes lesson by ID
// @Tags         admin
// @Produce      json
// @Param        id   path  int  true  "Lesson ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /lessons/{id} [delete]
func (h *Handler) DeleteLessonById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteLessonById(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
