package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/entity"

	"gorm.io/gorm"
)

func (s *Service) GetLessonById(ctx context.Context, id uint, userUUID string) (entity.Lesson, error) {

	lesson, err := s.repo.Lesson().GetLessonById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Lesson{}, fmt.Errorf("lesson with id %d not found", id)
		}
		return entity.Lesson{}, err
	}

	// Проверяем что урок реально найден
	if lesson == nil {
		return entity.Lesson{}, fmt.Errorf("lesson with id %d not found", id)
	}

	hasAccess, err := s.CheckUserAccessToLesson(ctx, userUUID, id)
	if err != nil {
		return entity.Lesson{}, fmt.Errorf("failed to check access: %w", err)
	}
	if !hasAccess {
		return entity.Lesson{}, fmt.Errorf("access denied")
	}

	return *lesson, nil
}
