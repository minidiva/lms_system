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

func (s *Service) CreateLesson(ctx context.Context, chapterId uint, lesson entity.Lesson) (uint, error) {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"chapter_id":  chapterId,
		"lesson_name": lesson.Name,
	}).Info("Creating new lesson")

	_, err := s.repo.Chapter().GetChapterById(ctx, chapterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("chapter_id", chapterId).Error("Chapter not found")
			return 0, fmt.Errorf("chapter with id %d not found", chapterId)
		}
		s.logger.WithError(err).Error("Failed to get chapter")
		return 0, err
	}

	lesson.ChapterID = chapterId
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()

	// DEBUG — детали создаваемого урока
	s.logger.WithFields(logrus.Fields{
		"chapter_id":     chapterId,
		"lesson_name":    lesson.Name,
		"order_position": lesson.OrderPosition,
	}).Debug("Lesson details")

	id, err := s.repo.Lesson().CreateLesson(ctx, chapterId, lesson)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create lesson")
		return 0, err
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"lesson_id":  id,
		"chapter_id": chapterId,
	}).Info("Lesson created successfully")

	return id, nil
}
