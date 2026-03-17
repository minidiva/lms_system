package chapter

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteChapterById godoc
// @Summary      Delete chapter
// @Description  Deletes chapter by ID
// @Tags         admin
// @Produce      json
// @Param        id   path  int  true  "Chapter ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /admin/chapter/delete/{id} [delete]
func (h *Handler) DeleteChapterById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chapter ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteChapterById(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
