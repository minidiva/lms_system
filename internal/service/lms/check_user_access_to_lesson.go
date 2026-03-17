package lms

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) CheckUserAccessToLesson(ctx context.Context, userId string, lessonId uint) (bool, error) {
	s.logger.WithFields(logrus.Fields{
		"user_id":   userId,
		"lesson_id": lessonId,
	}).Info("Checking user access to lesson")

	// Шаг 1 — получаем урок
	lesson, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get lesson")
		return false, fmt.Errorf("failed to get lesson: %w", err)
	}
	if lesson == nil {
		s.logger.WithField("lesson_id", lessonId).Error("Lesson not found")
		return false, fmt.Errorf("lesson not found")
	}

	s.logger.WithFields(logrus.Fields{
		"lesson_id":  lessonId,
		"chapter_id": lesson.ChapterID,
	}).Debug("Lesson details")

	// Шаг 2 — получаем главу
	chapter, err := s.repo.Chapter().GetChapterById(ctx, lesson.ChapterID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get chapter")
		return false, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		s.logger.WithField("chapter_id", lesson.ChapterID).Error("Chapter not found")
		return false, fmt.Errorf("chapter not found")
	}

	s.logger.WithFields(logrus.Fields{
		"chapter_id": lesson.ChapterID,
		"course_id":  chapter.CourseID,
	}).Debug("Chapter details")

	// Шаг 3 — проверяем доступ
	access, err := s.repo.UserCourseAccess().GetByUserIdAndCourseId(ctx, userId, chapter.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // нет доступа — не ошибка
		}
		s.logger.WithError(err).Error("Failed to check course access")
		return false, fmt.Errorf("failed to check course access: %w", err)
	}

	hasAccess := access != nil && access.Unlocked

	s.logger.WithFields(logrus.Fields{
		"user_id":    userId,
		"lesson_id":  lessonId,
		"has_access": hasAccess,
	}).Info("User access check completed")

	return hasAccess, nil
}
