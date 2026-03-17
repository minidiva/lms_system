package file

import (
	"lms_system/internal/domain"
)

type Handler struct {
	fileService domain.FileServiceInterface
}

func NewHandler(fileService domain.FileServiceInterface) *Handler {
	return &Handler{
		fileService: fileService,
	}
}
