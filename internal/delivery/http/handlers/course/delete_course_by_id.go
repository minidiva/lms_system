package course

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteCourseById godoc
// @Summary      Delete course
// @Description  Deletes course by ID
// @Tags         admin
// @Produce      json
// @Param        id   path  int  true  "Course ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /admin/courses/delete/{id} [delete]
func (h *Handler) DeleteCourseById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCourseById(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
