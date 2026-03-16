package lms

import (
	"context"
	"fmt"
	"lms_system/internal/domain/entity"

	"github.com/sirupsen/logrus"
)

func (s *Service) GetLessonAttachments(ctx context.Context, lessonId uint) ([]entity.Attachment, error) {

	// INFO — общее действие
	s.logger.WithField("lesson_id", lessonId).Info("Getting lesson attachments")

	lesson, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get lesson")
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}
	if lesson == nil {
		s.logger.WithField("lesson_id", lessonId).Error("Lesson not found")
		return nil, fmt.Errorf("lesson not found")
	}

	attachments, err := s.repo.Attachment().GetAttachmentsByLessonId(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get attachments")
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}

	// DEBUG — детали
	s.logger.WithFields(logrus.Fields{
		"lesson_id": lessonId,
		"count":     len(attachments),
	}).Debug("Attachments retrieved")

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"lesson_id": lessonId,
		"count":     len(attachments),
	}).Info("Lesson attachments retrieved successfully")

	return attachments, nil
}
