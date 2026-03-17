package chapter

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

func (r *Repository) CreateChapter(ctx context.Context, courseId uint, chapter *entity.Chapter) (uint, error) {
	chapter.CourseID = courseId

	err := r.db.WithContext(ctx).Create(chapter).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "CreateChapter",
			"course_id": courseId,
			"chapter":   chapter,
			"error":     err.Error(),
		}).Error("Failed to create chapter")
		return 0, err
	}

	r.logger.WithFields(logrus.Fields{
		"method":     "CreateChapter",
		"chapter_id": chapter.ID,
		"course_id":  courseId,
		"name":       chapter.Name,
	}).Info("Chapter created successfully")

	return chapter.ID, nil
}

func (r *Repository) UpdateChapterById(ctx context.Context, chapter *entity.Chapter) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Chapter{}).
		Where("id = ?", chapter.ID).
		Updates(chapter).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":     "UpdateChapterById",
			"chapter_id": chapter.ID,
			"error":      err.Error(),
		}).Error("Failed to update chapter")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":     "UpdateChapterById",
		"chapter_id": chapter.ID,
	}).Info("Chapter updated successfully")

	return nil
}

func (r *Repository) DeleteChapterById(ctx context.Context, chapterId uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", chapterId).
		Delete(&entity.Chapter{}).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":     "DeleteChapterById",
			"chapter_id": chapterId,
			"error":      err.Error(),
		}).Error("Failed to delete chapter")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":     "DeleteChapterById",
		"chapter_id": chapterId,
	}).Info("Chapter deleted successfully")

	return nil
}

func (r *Repository) GetChapterById(ctx context.Context, id uint) (*entity.Chapter, error) {
	var chapter entity.Chapter

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&chapter).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithFields(logrus.Fields{
				"method":     "GetChapterById",
				"chapter_id": id,
			}).Warn("Chapter not found")
		} else {
			r.logger.WithFields(logrus.Fields{
				"method":     "GetChapterById",
				"chapter_id": id,
				"error":      err.Error(),
			}).Error("Failed to get chapter")
		}
		return nil, err
	}

	return &chapter, nil
}

func (r *Repository) GetChaptersByCourseId(ctx context.Context, id uint) ([]entity.Chapter, error) {
	var chapters []entity.Chapter

	err := r.db.WithContext(ctx).
		Where("course_id = ?", id).
		Find(&chapters).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "GetChaptersByCourseId",
			"course_id": id,
			"error":     err.Error(),
		}).Error("Failed to get chapters by course ID")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "GetChaptersByCourseId",
		"course_id": id,
		"count":     len(chapters),
	}).Info("Successfully retrieved chapters by course ID")

	return chapters, nil
}
