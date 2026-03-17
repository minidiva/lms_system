package chapter

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lms_system/internal/domain/entity"

	"github.com/go-chi/chi/v5"
)

// CreateChapter godoc
// @Summary      Create chapter in course
// @Description  Creates a new chapter in the specified course
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        courseId  path      int             true  "Course ID"
// @Param        input     body      entity.Chapter  true  "Chapter data"
// @Success      201       {object}  map[string]uint
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Security     BearerAuth
// @Router       /courses/{courseId}/chapters [post]
func (h *Handler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	courseIdStr := chi.URLParam(r, "courseId")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var chapter entity.Chapter
	if err := json.NewDecoder(r.Body).Decode(&chapter); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateChapter(r.Context(), uint(courseId), chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]uint{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
