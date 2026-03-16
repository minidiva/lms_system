package lms

import (
	"context"
	"fmt"
	"lms_system/internal/utils"

	"github.com/sirupsen/logrus"
)

func (s *Service) CheckUserAccessToLesson(ctx context.Context, userId uint, lessonId uint) (bool, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"user_id":   userId,
		"lesson_id": lessonId,
	}).Info("Checking user access to lesson")

	if userId == 1 || userId == utils.ConvertKeycloakIDToUint("32bfb3d7-5b2c-4502-b08a-92ae81984f57") {
		s.logger.WithField("user_id", userId).Debug("Admin user access granted")
		return true, nil
	}

	lesson, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get lesson")
		return false, fmt.Errorf("failed to get lesson: %w", err)
	}
	if lesson == nil {
		s.logger.WithField("lesson_id", lessonId).Error("Lesson not found")
		return false, fmt.Errorf("lesson not found")
	}

	// DEBUG — детали урока
	s.logger.WithFields(logrus.Fields{
		"lesson_id":  lessonId,
		"chapter_id": lesson.ChapterID,
	}).Debug("Lesson details")

	chapter, err := s.repo.Chapter().GetChapterById(ctx, lesson.ChapterID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get chapter")
		return false, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		s.logger.WithField("chapter_id", lesson.ChapterID).Error("Chapter not found")
		return false, fmt.Errorf("chapter not found")
	}

	// DEBUG — детали главы
	s.logger.WithFields(logrus.Fields{
		"chapter_id": lesson.ChapterID,
		"course_id":  chapter.CourseID,
	}).Debug("Chapter details")

	access, err := s.repo.UserCourseAccess().GetByUserIdAndCourseId(ctx, userId, chapter.CourseID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to check course access")
		return false, fmt.Errorf("failed to check course access: %w", err)
	}

	hasAccess := access != nil && access.Unlocked

	// INFO — результат проверки
	s.logger.WithFields(logrus.Fields{
		"user_id":    userId,
		"lesson_id":  lessonId,
		"has_access": hasAccess,
	}).Info("User access check completed")

	return hasAccess, nil
}
