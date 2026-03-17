package course

import (
	"encoding/json"
	"lms_system/internal/domain/dto"
	"lms_system/utils"
	"net/http"
)

// BuyCourse godoc
// @Summary      Buy course
// @Description  Purchase a course for the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body      dto.BuyCourseRequest  true  "Course ID to buy"
// @Success      201    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/courses/buy [post]
func (h *Handler) BuyCourse(w http.ResponseWriter, r *http.Request) {
	userCtx := utils.GetUserFromContext(r.Context())
	if userCtx == nil || userCtx.UserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request dto.BuyCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// UUID берём из контекста, не из тела запроса
	request.UserUUID = userCtx.UserID

	if err := h.service.BuyCourse(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "Course purchased successfully"}
	json.NewEncoder(w).Encode(response)
}
