package file

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteFile handles file deletion
func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	// Get file path from request body
	var req struct {
		FilePath string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FilePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Delete file
	err := h.fileService.DeleteFile(r.Context(), req.FilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{
		"message": "File deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
