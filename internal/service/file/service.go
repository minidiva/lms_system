package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"lms_system/internal/domain"
	"lms_system/internal/infrastructure/clients/minio"
)

type Service struct {
	minioClient *minio.Client
}

func NewService(minioClient *minio.Client) domain.FileServiceInterface {
	return &Service{
		minioClient: minioClient,
	}
}

// UploadFile uploads a file and returns the file path/key
func (s *Service) UploadFile(ctx context.Context, fileType string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate unique object name
	objectName := s.generateObjectName(fileType, fileHeader.Filename)

	// Upload file to MinIO
	uploadedPath, err := s.minioClient.UploadFile(ctx, objectName, file, fileHeader)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadedPath, nil
}

// GetFile retrieves a file by its path/key
func (s *Service) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return s.minioClient.DownloadFile(ctx, filePath)
}

// DeleteFile removes a file by its path/key
func (s *Service) DeleteFile(ctx context.Context, filePath string) error {
	return s.minioClient.DeleteFile(ctx, filePath)
}

// GetFileURL returns a presigned URL for file access (valid for 1 hour)
func (s *Service) GetFileURL(ctx context.Context, filePath string) (string, error) {
	// URL valid for 1 hour
	expiry := 3600 // seconds
	return s.minioClient.GetFileURL(ctx, filePath, expiry)
}

// UploadLessonFile uploads a file associated with a lesson
func (s *Service) UploadLessonFile(ctx context.Context, lessonID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate object name with lesson prefix
	objectName := s.generateObjectName(fmt.Sprintf("lessons/%d", lessonID), fileHeader.Filename)

	// Upload file to MinIO
	uploadedPath, err := s.minioClient.UploadFile(ctx, objectName, file, fileHeader)
	if err != nil {
		return "", fmt.Errorf("failed to upload lesson file: %w", err)
	}

	return uploadedPath, nil
}

// UploadCourseFile uploads a file associated with a course
func (s *Service) UploadCourseFile(ctx context.Context, courseID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate object name with course prefix
	objectName := s.generateObjectName(fmt.Sprintf("courses/%d", courseID), fileHeader.Filename)

	// Upload file to MinIO
	uploadedPath, err := s.minioClient.UploadFile(ctx, objectName, file, fileHeader)
	if err != nil {
		return "", fmt.Errorf("failed to upload course file: %w", err)
	}

	return uploadedPath, nil
}

// generateObjectName generates a unique object name based on the original filename
func (s *Service) generateObjectName(prefix string, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)

	// Clean the name - remove special characters
	name = strings.ReplaceAll(name, " ", "_")

	// Add timestamp for uniqueness
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%s/%s_%d%s", prefix, name, timestamp, ext)
}
