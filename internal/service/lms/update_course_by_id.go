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

func (s *Service) UpdateCourseById(ctx context.Context, course entity.Course) error {

	// INFO — общее действие
	s.logger.WithField("course_id", course.ID).Info("Updating course")

	existingCourse, err := s.repo.Course().GetCourseById(ctx, course.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", course.ID).Error("Course not found")
			return fmt.Errorf("course with id %d not found", course.ID)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return err
	}

	course.UpdatedAt = time.Now()
	course.CreatedAt = existingCourse.Course.CreatedAt

	// DEBUG — детали обновления
	s.logger.WithFields(logrus.Fields{
		"course_id":   course.ID,
		"name":        course.Name,
		"description": course.Description,
	}).Debug("Course update details")

	err = s.repo.Course().UpdateCourseById(ctx, course)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update course")
		return err
	}

	// INFO — успешный результат
	s.logger.WithField("course_id", course.ID).Info("Course updated successfully")

	return nil
}
