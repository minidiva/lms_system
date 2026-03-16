package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) UpdateChapterById(ctx context.Context, chapter entity.Chapter) error {

	// INFO — общее действие
	s.logger.WithField("chapter_id", chapter.ID).Info("Updating chapter")

	existingChapter, err := s.repo.Chapter().GetChapterById(ctx, chapter.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("chapter_id", chapter.ID).Error("Chapter not found")
			return fmt.Errorf("chapter with id %d not found", chapter.ID)
		}
		s.logger.WithError(err).Error("Failed to get chapter")
		return err
	}

	chapter.UpdatedAt = time.Now()
	chapter.CreatedAt = existingChapter.CreatedAt
	chapter.CourseID = existingChapter.CourseID

	// DEBUG — детали обновления
	s.logger.WithFields(logrus.Fields{
		"chapter_id": chapter.ID,
		"name":       chapter.Name,
		"course_id":  chapter.CourseID,
	}).Debug("Chapter update details")

	err = s.repo.Chapter().UpdateChapterById(ctx, &chapter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update chapter")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("chapter_id", chapter.ID).Info("Chapter updated successfully")

	return nil
}
