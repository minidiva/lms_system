package lms

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (s *Service) DownloadAttachment(ctx context.Context, attachmentId uint, userId uint) (string, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"attachment_id": attachmentId,
		"user_id":       userId,
	}).Info("Downloading attachment")

	attachment, err := s.repo.Attachment().GetAttachmentById(ctx, attachmentId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get attachment")
		return "", fmt.Errorf("failed to get attachment: %w", err)
	}
	if attachment == nil {
		s.logger.WithField("attachment_id", attachmentId).Error("Attachment not found")
		return "", fmt.Errorf("attachment not found")
	}

	// DEBUG — детали вложения
	s.logger.WithFields(logrus.Fields{
		"attachment_id": attachmentId,
		"lesson_id":     attachment.LessonID,
		"url":           attachment.URL,
	}).Debug("Attachment details")

	hasAccess, err := s.CheckUserAccessToLesson(ctx, userId, attachment.LessonID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to check user access")
		return "", fmt.Errorf("failed to check access: %w", err)
	}
	if !hasAccess {
		s.logger.WithFields(logrus.Fields{
			"user_id":       userId,
			"attachment_id": attachmentId,
		}).Error("Access denied")
		return "", fmt.Errorf("access denied")
	}

	url, err := s.fileService.GetFileURL(ctx, attachment.URL)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get download URL")
		return "", fmt.Errorf("failed to get download URL: %w", err)
	}

	// INFO — успешный результат
	s.logger.WithField("attachment_id", attachmentId).Info("Attachment download URL generated successfully")

	return url, nil
}
