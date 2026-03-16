package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/entity"

	"gorm.io/gorm"
)

func (s *Service) GetCourseById(ctx context.Context, id uint) (entity.CourseAggregate, error) {

	// INFO — общее действие
	s.logger.WithField("course_id", id).Info("Getting course by id")

	courseAggregate, err := s.repo.Course().GetCourseById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", id).Error("Course not found")
			return entity.CourseAggregate{}, fmt.Errorf("course with id %d not found", id)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return entity.CourseAggregate{}, err
	}

	// DEBUG — детали курса
	s.logger.WithField("course_id", id).Debug("Course details retrieved")

	// INFO — успешный результат
	s.logger.WithField("course_id", id).Info("Course retrieved successfully")

	return *courseAggregate, nil
}
