package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/entity"

	"gorm.io/gorm"
)

func (s *Service) GetLessonById(ctx context.Context, id uint) (entity.Lesson, error) {

	// INFO — общее действие
	s.logger.WithField("lesson_id", id).Info("Getting lesson by id")

	lesson, err := s.repo.Lesson().GetLessonById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("lesson_id", id).Error("Lesson not found")
			return entity.Lesson{}, fmt.Errorf("lesson with id %d not found", id)
		}
		s.logger.WithError(err).Error("Failed to get lesson")
		return entity.Lesson{}, err
	}

	// DEBUG — детали урока
	s.logger.WithField("lesson_id", id).Debug("Lesson details retrieved")

	// INFO — успешный результат
	s.logger.WithField("lesson_id", id).Info("Lesson retrieved successfully")

	return *lesson, nil
}
