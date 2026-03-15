package course

import (
	"context"
	"errors"
	"lms_system/internal/domain/entity"

	"github.com/sirupsen/logrus"
	_ "gorm.io/driver/postgres"
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

func (r *Repository) CreateCourse(ctx context.Context, course entity.Course) (uint, error) {
	if err := r.db.WithContext(ctx).Create(&course).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"method": "CreateCourse",
			"course": course,
			"error":  err.Error(),
		}).Error("Failed to create course")
		return 0, err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "CreateCourse",
		"course_id": course.ID,
		"name":      course.Name,
	}).Info("Course created successfully")

	return course.ID, nil
}

func (r *Repository) UpdateCourseById(ctx context.Context, course entity.Course) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Course{}).
		Where("id = ?", course.ID).
		Updates(course).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "UpdateCourseById",
			"course_id": course.ID,
			"error":     err.Error(),
		}).Error("Failed to update course")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "UpdateCourseById",
		"course_id": course.ID,
	}).Info("Course updated successfully")

	return nil
}

func (r *Repository) DeleteCourseById(ctx context.Context, courseId uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", courseId).
		Delete(&entity.Course{}).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "DeleteCourseById",
			"course_id": courseId,
			"error":     err.Error(),
		}).Error("Failed to delete course")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"method":    "DeleteCourseById",
		"course_id": courseId,
	}).Info("Course deleted successfully")

	return nil
}

func (r *Repository) GetCourseById(ctx context.Context, id uint) (*entity.CourseAggregate, error) {

	// Шаг 1: загружаем курс
	var course entity.Course
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&course).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetCourseById",
				"course_id": id,
			}).Warn("Course not found")
		} else {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetCourseById",
				"course_id": id,
				"error":     err.Error(),
			}).Error("Failed to get course")
		}
		return nil, err
	}

	// Шаг 2: загружаем главы этого курса отдельно
	var chapters []entity.Chapter
	r.db.WithContext(ctx).
		Where("course_id = ?", id).
		Find(&chapters)

	// Шаг 3: собираем агрегат вручную
	return &entity.CourseAggregate{
		Course:   course,
		Chapters: chapters,
	}, nil
}

func (r *Repository) GetAllCourses(ctx context.Context) ([]entity.Course, error) {
	var courses []entity.Course

	err := r.db.WithContext(ctx).
		Find(&courses).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method": "GetAllCourses",
			"error":  err.Error(),
		}).Error("Failed to get all courses")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"method": "GetAllCourses",
		"count":  len(courses),
	}).Info("Successfully retrieved all courses")

	return courses, nil
}
