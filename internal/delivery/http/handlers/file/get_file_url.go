package file

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetFileURL returns a presigned URL for file access
func (h *Handler) GetFileURL(w http.ResponseWriter, r *http.Request) {
	// Get file path from query parameter
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Get presigned URL
	url, err := h.fileService.GetFileURL(r.Context(), filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{
		"url":       url,
		"file_path": filePath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
