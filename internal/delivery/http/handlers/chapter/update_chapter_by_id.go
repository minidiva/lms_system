package chapter

import (
	"encoding/json"
	"net/http"
	"strconv"

	dto "lms_system/internal/domain/dto/chapter"
	"lms_system/internal/domain/entity"

	"github.com/go-chi/chi/v5"
)

// UpdateChapterById godoc
// @Summary      Update chapter
// @Description  Updates chapter by ID
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Chapter ID"
// @Param        input  body      entity.Chapter  true  "Chapter data"
// @Success      204
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /chapters/update/{id} [put]
func (h *Handler) UpdateChapterById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chapter ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chapter := entity.Chapter{
		Name:          req.Name,
		Description:   req.Description,
		OrderPosition: req.OrderPosition,
	}

	chapter.ID = uint(id)
	if err := h.service.UpdateChapterById(r.Context(), chapter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
