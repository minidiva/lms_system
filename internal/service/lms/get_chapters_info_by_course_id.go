package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) GetChaptersInfoByCourseId(ctx context.Context, id uint) ([]entity.ChapterInfoAggregate, error) {

	// INFO — общее действие
	s.logger.WithField("course_id", id).Info("Getting chapters info by course id")

	_, err := s.repo.Course().GetCourseById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", id).Error("Course not found")
			return nil, fmt.Errorf("course with id %d not found", id)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return nil, err
	}

	chapters, err := s.repo.Chapter().GetChaptersByCourseId(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get chapters")
		return nil, err
	}

	// DEBUG — детали
	s.logger.WithFields(logrus.Fields{
		"course_id":      id,
		"chapters_count": len(chapters),
	}).Debug("Chapters retrieved")

	chaptersInfo := make([]entity.ChapterInfoAggregate, 0, len(chapters))
	for _, chapter := range chapters {
		lessons, err := s.repo.Lesson().GetAllLessonsByChapterId(ctx, chapter.ID)
		if err != nil {
			s.logger.WithError(err).Error("Failed to get lessons for chapter")
			return nil, err
		}

		lessonNames := make([]string, 0, len(lessons))
		for _, lesson := range lessons {
			lessonNames = append(lessonNames, lesson.Name)
		}

		chapterInfo := entity.ChapterInfoAggregate{
			Chapter:     chapter,
			LessonsName: lessonNames,
		}
		chaptersInfo = append(chaptersInfo, chapterInfo)
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"course_id": id,
		"count":     len(chaptersInfo),
	}).Info("Chapters info retrieved successfully")

	return chaptersInfo, nil
}
