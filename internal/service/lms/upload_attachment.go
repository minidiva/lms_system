package lms

import (
	"context"
	"fmt"
	"lms_system/internal/domain/entity"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

func (s *Service) UploadAttachment(ctx context.Context, lessonId uint, file multipart.File, fileHeader *multipart.FileHeader) (*entity.Attachment, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"lesson_id": lessonId,
		"filename":  fileHeader.Filename,
	}).Info("Uploading attachment to lesson")

	lesson, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get lesson")
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}
	if lesson == nil {
		s.logger.WithField("lesson_id", lessonId).Error("Lesson not found")
		return nil, fmt.Errorf("lesson not found")
	}

	filePath, err := s.fileService.UploadLessonFile(ctx, lessonId, file, fileHeader)
	if err != nil {
		s.logger.WithError(err).Error("Failed to upload file to storage")
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	attachment := &entity.Attachment{
		Name:     fileHeader.Filename,
		URL:      filePath,
		LessonID: lessonId,
	}

	// DEBUG — детали вложения
	s.logger.WithFields(logrus.Fields{
		"lesson_id": lessonId,
		"filename":  fileHeader.Filename,
		"file_path": filePath,
	}).Debug("Attachment details")

	if err := s.repo.Attachment().CreateAttachment(ctx, attachment); err != nil {
		s.logger.WithError(err).Error("Failed to save attachment to database")
		_ = s.fileService.DeleteFile(ctx, filePath)
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"lesson_id": lessonId,
		"filename":  fileHeader.Filename,
	}).Info("Attachment uploaded successfully")

	return attachment, nil
}
