package course

import (
	"encoding/json"
	dto "lms_system/internal/domain/dto/course"
	"lms_system/internal/domain/entity"
	"net/http"
)

// CreateCourse godoc
// @Summary      Create course
// @Description  Creates a new course
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateCourseRequest  true  "Course data"
// @Success      201    {object}  map[string]uint
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /courses/create [post]
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	course := entity.Course{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := h.service.CreateCourse(r.Context(), course)
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
