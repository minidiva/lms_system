package lms

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (s *Service) DeleteAttachment(ctx context.Context, attachmentId uint) error {

	// INFO — общее действие
	s.logger.WithField("attachment_id", attachmentId).Info("Deleting attachment")

	attachment, err := s.repo.Attachment().GetAttachmentById(ctx, attachmentId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get attachment")
		return fmt.Errorf("failed to get attachment: %w", err)
	}
	if attachment == nil {
		s.logger.WithField("attachment_id", attachmentId).Error("Attachment not found")
		return fmt.Errorf("attachment not found")
	}

	// DEBUG — детали удаляемого вложения
	s.logger.WithFields(logrus.Fields{
		"attachment_id": attachmentId,
		"name":          attachment.Name,
		"url":           attachment.URL,
	}).Debug("Attachment details")

	// Delete from database
	if err := s.repo.Attachment().DeleteAttachment(ctx, attachmentId); err != nil {
		s.logger.WithError(err).Error("Failed to delete attachment from database")
		return fmt.Errorf("failed to delete attachment record: %w", err)
	}

	// Delete file from MinIO
	if err := s.fileService.DeleteFile(ctx, attachment.URL); err != nil {
		s.logger.WithError(err).Errorf("Failed to delete file from storage: %s", attachment.URL)
	}

	// INFO — успешный результат
	s.logger.WithField("attachment_id", attachmentId).Info("Attachment deleted successfully")

	return nil
}
