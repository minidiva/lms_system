package lms

import (
	"context"
	"fmt"
	"lms_system/internal/domain/entity"

	"github.com/sirupsen/logrus"
)

func (s *Service) GetAttachment(ctx context.Context, attachmentId uint, userId string) (*entity.Attachment, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"attachment_id": attachmentId,
		"user_id":       userId,
	}).Info("Getting attachment")

	attachment, err := s.repo.Attachment().GetAttachmentById(ctx, attachmentId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get attachment")
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	if attachment == nil {
		s.logger.WithField("attachment_id", attachmentId).Error("Attachment not found")
		return nil, fmt.Errorf("attachment not found")
	}

	// DEBUG — детали вложения
	s.logger.WithFields(logrus.Fields{
		"attachment_id": attachmentId,
		"name":          attachment.Name,
		"lesson_id":     attachment.LessonID,
	}).Debug("Attachment details")

	hasAccess, err := s.CheckUserAccessToLesson(ctx, userId, attachment.LessonID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to check user access")
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !hasAccess {
		s.logger.WithFields(logrus.Fields{
			"user_id":       userId,
			"attachment_id": attachmentId,
		}).Error("Access denied")
		return nil, fmt.Errorf("access denied")
	}

	// INFO — успешный результат
	s.logger.WithField("attachment_id", attachmentId).Info("Attachment retrieved successfully")

	return attachment, nil
}
