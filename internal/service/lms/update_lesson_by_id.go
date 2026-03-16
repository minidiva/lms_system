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

func (s *Service) UpdateLessonById(ctx context.Context, lesson entity.Lesson) error {

	// INFO — общее действие
	s.logger.WithField("lesson_id", lesson.ID).Info("Updating lesson")

	existingLesson, err := s.repo.Lesson().GetLessonById(ctx, lesson.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("lesson_id", lesson.ID).Error("Lesson not found")
			return fmt.Errorf("lesson with id %d not found", lesson.ID)
		}
		s.logger.WithError(err).Error("Failed to get lesson")
		return err
	}

	lesson.UpdatedAt = time.Now()
	lesson.CreatedAt = existingLesson.CreatedAt
	lesson.ChapterID = existingLesson.ChapterID

	// DEBUG — детали обновления
	s.logger.WithFields(logrus.Fields{
		"lesson_id":  lesson.ID,
		"name":       lesson.Name,
		"chapter_id": lesson.ChapterID,
	}).Debug("Lesson update details")

	err = s.repo.Lesson().UpdateLessonById(ctx, lesson)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update lesson")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("lesson_id", lesson.ID).Info("Lesson updated successfully")

	return nil
}
