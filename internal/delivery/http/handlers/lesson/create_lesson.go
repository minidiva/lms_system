package lesson

import (
	"encoding/json"
	"lms_system/internal/domain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CreateLesson godoc
// @Summary      Create lesson in chapter
// @Description  Creates a new lesson in the specified chapter
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        chapterId  path      int            true  "Chapter ID"
// @Param        input      body      entity.Lesson  true  "Lesson data"
// @Success      201        {object}  map[string]uint
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Security     BearerAuth
// @Router       /chapters/{chapterId}/lessons [post]
func (h *Handler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	chapterIdStr := chi.URLParam(r, "chapterId")
	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chapter ID", http.StatusBadRequest)
		return
	}

	var lesson entity.Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateLesson(r.Context(), uint(chapterId), lesson)
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
