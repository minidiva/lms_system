package lms

import (
	"context"
	"lms_system/internal/domain/entity"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Service) CreateCourse(ctx context.Context, course entity.Course) (uint, error) {

	// INFO — общее действие
	s.logger.WithField("course_name", course.Name).Info("Creating new course")

	// DEBUG — детали создаваемого курса
	s.logger.WithFields(logrus.Fields{
		"name":        course.Name,
		"description": course.Description,
	}).Debug("Course details")

	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()

	id, err := s.repo.Course().CreateCourse(ctx, course)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create course")
		return 0, err
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"course_id":   id,
		"course_name": course.Name,
	}).Info("Course created successfully")

	return id, nil
}
