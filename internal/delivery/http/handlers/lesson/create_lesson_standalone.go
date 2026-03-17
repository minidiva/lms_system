package lesson

import (
	"encoding/json"
	"lms_system/internal/domain/entity"
	"net/http"
)

type CreateLessonRequest struct {
	ChapterID     uint   `json:"chapter_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	OrderPosition int    `json:"order_position"`
}

// CreateLessonStandalone godoc
// @Summary      Create lesson standalone
// @Description  Creates a new lesson with chapter_id in request body
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        input  body      CreateLessonRequest    true  "Lesson data"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /lessons [post]
func (h *Handler) CreateLessonStandalone(w http.ResponseWriter, r *http.Request) {
	var req CreateLessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ChapterID == 0 {
		http.Error(w, "chapter_id is required", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	lesson := entity.Lesson{
		Title:         req.Title,
		Description:   req.Description,
		Content:       req.Content,
		OrderPosition: req.OrderPosition,
	}

	id, err := h.service.CreateLesson(r.Context(), req.ChapterID, lesson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         id,
		"chapter_id": req.ChapterID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
