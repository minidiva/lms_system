package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetChaptersInfoByCourseId godoc
// @Summary      Get chapters by course ID
// @Description  Returns list of chapters with lesson names for the specified course
// @Tags         courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {array}   entity.ChapterInfoAggregate
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /public/courses/{id}/chapters [get]
func (h *Handler) GetChaptersInfoByCourseId(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	chapters, err := h.service.GetChaptersInfoByCourseId(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chapters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
