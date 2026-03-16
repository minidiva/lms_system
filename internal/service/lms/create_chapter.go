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

func (s *Service) CreateChapter(ctx context.Context, courseId uint, chapter entity.Chapter) (uint, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"course_id":    courseId,
		"chapter_name": chapter.Name,
	}).Info("Creating new chapter")

	_, err := s.repo.Course().GetCourseById(ctx, courseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", courseId).Error("Course not found")
			return 0, fmt.Errorf("course with id %d not found", courseId)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return 0, err
	}

	chapter.CourseID = courseId
	chapter.CreatedAt = time.Now()
	chapter.UpdatedAt = time.Now()

	// DEBUG — детали создаваемой главы
	s.logger.WithFields(logrus.Fields{
		"course_id":      courseId,
		"chapter_name":   chapter.Name,
		"order_position": chapter.OrderPosition,
	}).Debug("Chapter details")

	id, err := s.repo.Chapter().CreateChapter(ctx, courseId, &chapter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create chapter")
		return 0, err
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"chapter_id": id,
		"course_id":  courseId,
	}).Info("Chapter created successfully")

	return id, nil
}
