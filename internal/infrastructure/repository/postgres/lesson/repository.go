package lesson

import (
	"context"
	"errors"
	"lms_system/internal/domain/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRepository(db *gorm.DB, logger *logrus.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateLesson(ctx context.Context, chapterId uint, lesson entity.Lesson) (uint, error) {
	lesson.ChapterID = chapterId

	err := r.db.WithContext(ctx).Create(&lesson).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":     "CreateLesson",
			"chapter_id": chapterId,
			"lesson":     lesson,
			"error":      err.Error(),
		}).Error("Failed to create lesson")
		return 0, err
	}

	r.logger.WithFields(logrus.Fields{
		"method":     "CreateLesson",
		"lesson_id":  lesson.ID,
		"chapter_id": chapterId,
	}).Info("Lesson created successfully")

	return lesson.ID, nil
}

func (r *Repository) UpdateLessonById(ctx context.Context, lesson entity.Lesson) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Lesson{}).
		Where("id = ?", lesson.ID).
		Updates(lesson).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "UpdateLessonById",
			"lesson_id": lesson.ID,
			"error":     err.Error(),
		}).Error("Failed to update lesson")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "UpdateLessonById",
		"lesson_id": lesson.ID,
	}).Info("Lesson updated successfully")

	return nil
}

func (r *Repository) DeleteLessonById(ctx context.Context, lessonId uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", lessonId).
		Delete(&entity.Lesson{}).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "DeleteLessonById",
			"lesson_id": lessonId,
			"error":     err.Error(),
		}).Error("Failed to delete lesson")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "DeleteLessonById",
		"lesson_id": lessonId,
	}).Info("Lesson deleted successfully")

	return nil
}

func (r *Repository) GetLessonById(ctx context.Context, id uint) (*entity.Lesson, error) {
	var lesson entity.Lesson

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&lesson).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetLessonById",
				"lesson_id": id,
			}).Warn("Lesson not found")
		} else {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetLessonById",
				"lesson_id": id,
				"error":     err.Error(),
			}).Error("Failed to get lesson")
		}
		return nil, err
	}

	return &lesson, nil
}

func (r *Repository) GetAllLessonsByChapterId(ctx context.Context, chapterId uint) ([]entity.Lesson, error) {
	var lessons []entity.Lesson

	err := r.db.WithContext(ctx).
		Where("chapter_id = ?", chapterId).
		Find(&lessons).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":     "GetAllLessonsByChapterId",
			"chapter_id": chapterId,
			"error":      err.Error(),
		}).Error("Failed to get lessons by chapter ID")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"method":     "GetAllLessonsByChapterId",
		"chapter_id": chapterId,
		"count":      len(lessons),
	}).Info("Successfully retrieved lessons by chapter ID")

	return lessons, nil
}
