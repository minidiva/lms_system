package lms

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (s *Service) DeleteLessonById(ctx context.Context, lessonId uint) error {

	// INFO — общее действие
	s.logger.WithField("lesson_id", lessonId).Info("Deleting lesson")

	_, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("lesson_id", lessonId).Error("Lesson not found")
			return fmt.Errorf("lesson with id %d not found", lessonId)
		}
		s.logger.WithError(err).Error("Failed to get lesson")
		return err
	}

	// DEBUG
	s.logger.WithField("lesson_id", lessonId).Debug("Deleting lesson from database")

	err = s.repo.Lesson().DeleteLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete lesson")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("lesson_id", lessonId).Info("Lesson deleted successfully")

	return nil
}
