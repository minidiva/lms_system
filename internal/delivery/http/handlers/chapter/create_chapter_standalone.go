package chapter

import (
	"encoding/json"
	dto "lms_system/internal/domain/dto/chapter"
	"lms_system/internal/domain/entity"
	"net/http"
)

// CreateChapterStandalone godoc
// @Summary      Create chapter standalone
// @Description  Creates a new chapter with course_id in request body
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        input  body      CreateChapterRequest  true  "Chapter data"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /chapters/create [post]
func (h *Handler) CreateChapterStandalone(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CourseID == 0 {
		http.Error(w, "course_id is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	chapter := entity.Chapter{
		Name:          req.Name,
		Description:   req.Description,
		OrderPosition: req.OrderPosition,
	}

	id, err := h.service.CreateChapter(r.Context(), req.CourseID, chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":        id,
		"course_id": req.CourseID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
