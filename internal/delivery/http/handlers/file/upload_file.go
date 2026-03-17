package file

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UploadFile handles generic file upload
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get file type from form (optional, defaults to "general")
	fileType := r.FormValue("type")
	if fileType == "" {
		fileType = "general"
	}

	// Upload file
	filePath, err := h.fileService.UploadFile(r.Context(), fileType, file, fileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{
		"file_path": filePath,
		"message":   "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
