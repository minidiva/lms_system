package user_course_access

import (
	"context"
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

func (r *Repository) GetAllByUserId(ctx context.Context, userId string) ([]entity.UserCourseAccess, error) {
	var accessList []entity.UserCourseAccess
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&accessList).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":  "GetAllByUserId",
			"user_id": userId,
			"error":   err.Error(),
		}).Error("Failed to get user course access by user ID")
		return nil, err
	}
	r.logger.WithFields(logrus.Fields{
		"method":  "GetAllByUserId",
		"user_id": userId,
		"count":   len(accessList),
	}).Info("Successfully retrieved user course access by user ID")
	return accessList, nil
}

func (r *Repository) UpdateAccess(ctx context.Context, access *entity.UserCourseAccess) error {
	err := r.db.WithContext(ctx).
		Model(&entity.UserCourseAccess{}).
		Where("user_id = ? AND course_id = ?", access.UserID, access.CourseID).
		Updates(access).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":    "UpdateAccess",
			"user_id":   access.UserID,
			"course_id": access.CourseID,
			"error":     err.Error(),
		}).Error("Failed to update user course access")
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"method":    "UpdateAccess",
		"user_id":   access.UserID,
		"course_id": access.CourseID,
	}).Info("User course access updated successfully")
	return nil
}

func (r *Repository) GetByUserIdAndCourseId(ctx context.Context, userId string, courseId uint) (*entity.UserCourseAccess, error) {
	var access entity.UserCourseAccess
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userId, courseId).
		First(&access).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetByUserIdAndCourseId",
				"user_id":   userId,
				"course_id": courseId,
			}).Warn("User course access not found")
		} else {
			r.logger.WithFields(logrus.Fields{
				"method":    "GetByUserIdAndCourseId",
				"user_id":   userId,
				"course_id": courseId,
				"error":     err.Error(),
			}).Error("Failed to get user course access")
		}
		return nil, err
	}
	return &access, nil
}

func (r *Repository) CreateUserCourseAccess(ctx context.Context, userAccess entity.UserCourseAccess) error {
	err := r.db.WithContext(ctx).Create(&userAccess).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"method":      "CreateUserCourseAccess",
			"user_id":     userAccess.UserID,
			"course_id":   userAccess.CourseID,
			"user_access": userAccess,
			"error":       err.Error(),
		}).Error("Failed to create user course access")
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"method":    "CreateUserCourseAccess",
		"user_id":   userAccess.UserID,
		"course_id": userAccess.CourseID,
	}).Info("User course access created successfully")
	return nil
}
