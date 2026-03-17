package lesson

import (
	"encoding/json"
	"lms_system/internal/domain/dto"
	"lms_system/internal/domain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// UpdateLessonById godoc
// @Summary      Update lesson
// @Description  Updates lesson by ID
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        id     path      int            true  "Lesson ID"
// @Param        input  body      dto.UpdateLessonRequest  true  "Lesson data"
// @Success      204
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /lessons/{id} [put]
func (h *Handler) UpdateLessonById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateLessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lesson := entity.Lesson{
		ID:            uint(id),
		Name:          req.Name,
		Description:   req.Description,
		Content:       req.Content,
		OrderPosition: req.OrderPosition,
	}

	if err := h.service.UpdateLessonById(r.Context(), lesson); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
