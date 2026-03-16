package lms

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (s *Service) DeleteCourseById(ctx context.Context, courseId uint) error {

	// INFO — общее действие
	s.logger.WithField("course_id", courseId).Info("Deleting course")

	_, err := s.repo.Course().GetCourseById(ctx, courseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", courseId).Error("Course not found")
			return fmt.Errorf("course with id %d not found", courseId)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return err
	}

	// DEBUG
	s.logger.WithField("course_id", courseId).Debug("Deleting course from database")

	err = s.repo.Course().DeleteCourseById(ctx, courseId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete course")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("course_id", courseId).Info("Course deleted successfully")

	return nil
}
