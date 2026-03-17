package lms

import (
	"context"
	"errors"
	"fmt"
	"lms_system/internal/domain/dto"
	"lms_system/internal/domain/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) BuyCourse(ctx context.Context, request dto.BuyCourseRequest) error {

	// INFO — общее действие
	s.logger.WithFields(logrus.Fields{
		"user_id":   request.UserUUID,
		"course_id": request.CourseId,
	}).Info("User buying course")

	_, err := s.repo.Course().GetCourseById(ctx, request.CourseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithField("course_id", request.CourseId).Error("Course not found")
			return fmt.Errorf("course with id %d not found", request.CourseId)
		}
		s.logger.WithError(err).Error("Failed to get course")
		return err
	}

	existingAccess, err := s.repo.UserCourseAccess().GetByUserIdAndCourseId(ctx, request.UserUUID, request.CourseId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.WithError(err).Error("Failed to check existing course access")
		return err
	}

	if existingAccess != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":   request.UserUUID,
			"course_id": request.CourseId,
		}).Error("User already has access to this course")
		return fmt.Errorf("user already has access to this course")
	}

	userAccess := entity.UserCourseAccess{
		UserID:    request.UserUUID,
		CourseID:  request.CourseId,
		Unlocked:  true,
		CreatedAt: time.Now(),
	}

	// DEBUG — детали создаваемого доступа
	s.logger.WithFields(logrus.Fields{
		"user_id":   userAccess.UserID,
		"course_id": userAccess.CourseID,
		"unlocked":  userAccess.Unlocked,
	}).Debug("Creating user course access")

	err = s.repo.UserCourseAccess().CreateUserCourseAccess(ctx, userAccess)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create user course access")
		return err
	}

	// INFO — успешный результат
	s.logger.WithFields(logrus.Fields{
		"user_id":   request.UserUUID,
		"course_id": request.CourseId,
	}).Info("Course purchased successfully")

	return nil
}
