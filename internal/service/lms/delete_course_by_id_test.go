package lms

import (
	"context"
	"errors"
	"lms_system/internal/domain/entity"
	"lms_system/internal/service/lms/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeleteCourseById_Success(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)
	mockCourseRepo.On("DeleteCourseById", mock.Anything, uint(1)).
		Return(nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteCourseById(context.Background(), uint(1))

	assert.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseById_NotFound(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteCourseById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockCourseRepo.AssertNotCalled(t, "DeleteCourseById")
}

func TestDeleteCourseById_GetCourseDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	dbError := errors.New("connection refused")
	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteCourseById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	mockCourseRepo.AssertNotCalled(t, "DeleteCourseById")
}

func TestDeleteCourseById_DeleteDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)

	dbError := errors.New("delete failed")
	mockCourseRepo.On("DeleteCourseById", mock.Anything, uint(1)).
		Return(dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteCourseById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
}
