package attachment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lms_system/internal/domain"
	"lms_system/internal/domain/common"
	"lms_system/internal/domain/entity"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service domain.ServiceInterface
}

func NewHandler(service domain.ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := chi.URLParam(r, "lessonId")
	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}
	err = r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	attachment, err := h.service.UploadAttachment(r.Context(), uint(lessonID), file, fileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload attachment: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attachment)
}

func (h *Handler) DownloadAttachment(w http.ResponseWriter, r *http.Request) {
	attachmentIDStr := chi.URLParam(r, "attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid attachment ID", http.StatusBadRequest)
		return
	}
	userCtx, ok := r.Context().Value(common.UserContextKey).(*entity.UserContext)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userCtx.UserID
	downloadURL, err := h.service.DownloadAttachment(r.Context(), uint(attachmentID), userID)
	if err != nil {
		if err.Error() == "access denied" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		if err.Error() == "attachment not found" {
			http.Error(w, "Attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get download URL: %v", err), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, downloadURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GetLessonAttachments(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := chi.URLParam(r, "lessonId")
	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}
	attachments, err := h.service.GetLessonAttachments(r.Context(), uint(lessonID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get attachments: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachments)
}

func (h *Handler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	attachmentIDStr := chi.URLParam(r, "attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid attachment ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteAttachment(r.Context(), uint(attachmentID))
	if err != nil {
		if err.Error() == "attachment not found" {
			http.Error(w, "Attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete attachment: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Attachment deleted successfully",
	})
}
