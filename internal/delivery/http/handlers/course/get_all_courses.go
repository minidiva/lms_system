package course

import (
	"encoding/json"
	"net/http"
)

// GetAllCourses godoc
// @Summary      Get all courses
// @Description  Returns list of all courses
// @Tags         courses
// @Produce      json
// @Success      200  {array}   entity.Course
// @Failure      500  {object}  map[string]string
// @Router       /public/courses [get]
func (h *Handler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.GetAllCourses(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
