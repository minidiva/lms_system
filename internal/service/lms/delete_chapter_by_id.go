package lms

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) DeleteChapterById(ctx context.Context, chapterId uint) error {

	// INFO — общее действие
	s.logger.WithField("chapter_id", chapterId).Info("Deleting chapter")

	_, err := s.repo.Chapter().GetChapterById(ctx, chapterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("chapter_id", chapterId).Error("Chapter not found")
			return fmt.Errorf("chapter with id %d not found", chapterId)
		}
		s.logger.WithError(err).Error("Failed to get chapter")
		return err
	}

	// DEBUG — детали удаляемой главы
	s.logger.WithFields(logrus.Fields{
		"chapter_id": chapterId,
	}).Debug("Deleting chapter from database")

	err = s.repo.Chapter().DeleteChapterById(ctx, chapterId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete chapter")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("chapter_id", chapterId).Info("Chapter deleted successfully")

	return nil
}
