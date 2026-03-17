package file

import (
	"fmt"
	"io"
	"net/http"
)

// DownloadFile handles file download
func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Get file path from query parameter
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Get file from storage
	fileReader, err := h.fileService.GetFile(r.Context(), filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusNotFound)
		return
	}
	defer fileReader.Close()

	// Set appropriate headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filePath))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy file to response
	_, err = io.Copy(w, fileReader)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
