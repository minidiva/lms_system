package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/common"
	"lms_system/internal/domain/entity"

	"gorm.io/gorm"
)

func (s *Service) GetLessonById(ctx context.Context, id uint, userUUID string, role common.UserRole) (entity.Lesson, error) {

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

	if lesson == nil {
		return entity.Lesson{}, fmt.Errorf("lesson with id %d not found", id)
	}

	hasAccess, err := s.CheckUserAccessToLesson(ctx, userUUID, role, id)
	if err != nil {
		return entity.Lesson{}, fmt.Errorf("failed to check access: %w", err)
	}
	if !hasAccess {
		return entity.Lesson{}, fmt.Errorf("access denied")
	}

	s.logger.WithField("lesson_id", id).Info("Lesson retrieved successfully")

	return *lesson, nil
}
