package course

import (
	"encoding/json"
	dto "lms_system/internal/domain/dto/course"
	"lms_system/internal/domain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// UpdateCourseById godoc
// @Summary      Update course
// @Description  Updates course by ID
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        id     path      int                      true  "Course ID"
// @Param        input  body      dto.CreateCourseRequest  true  "Course data"
// @Success      204
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /courses/update/{id} [put]
func (h *Handler) UpdateCourseById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	course := entity.Course{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.UpdateCourseById(r.Context(), course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
