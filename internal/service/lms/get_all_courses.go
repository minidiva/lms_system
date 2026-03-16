package lms

import (
	"context"
	"lms_system/internal/domain/entity"
)

func (s *Service) GetAllCourses(ctx context.Context) ([]entity.Course, error) {

	// INFO — общее действие
	s.logger.Info("Getting all courses")

	courses, err := s.repo.Course().GetAllCourses(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get all courses")
		return nil, err
	}

	// DEBUG — детали результата
	s.logger.WithField("count", len(courses)).Debug("Courses retrieved")

	// INFO — успешный результат
	s.logger.WithField("count", len(courses)).Info("All courses retrieved successfully")

	return courses, nil
}
