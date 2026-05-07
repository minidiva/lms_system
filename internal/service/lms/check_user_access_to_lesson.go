package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/common"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) CheckUserAccessToLesson(ctx context.Context, userId string, role common.UserRole, lessonId uint) (bool, error) {

	s.logger.WithFields(logrus.Fields{
		"user_id":   userId,
		"lesson_id": lessonId,
		"role":      role,
	}).Info("Checking user access to lesson")

	// Админ и учитель имеют доступ ко всем урокам
	if role == common.RoleAdmin || role == common.RoleTeacher {
		s.logger.WithField("user_id", userId).Debug("Admin/Teacher access granted")
		return true, nil
	}

	// Остальная логика для обычных пользователей
	lesson, err := s.repo.Lesson().GetLessonById(ctx, lessonId)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get lesson")
		return false, fmt.Errorf("failed to get lesson: %w", err)
	}
	if lesson == nil {
		return false, fmt.Errorf("lesson not found")
	}

	chapter, err := s.repo.Chapter().GetChapterById(ctx, lesson.ChapterID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get chapter")
		return false, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return false, fmt.Errorf("chapter not found")
	}

	access, err := s.repo.UserCourseAccess().GetByUserIdAndCourseId(ctx, userId, chapter.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
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
