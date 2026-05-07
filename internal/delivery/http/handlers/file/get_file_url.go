package file

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetFileURL(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	url, err := h.fileService.GetFileURL(r.Context(), filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file URL: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
	})
}
