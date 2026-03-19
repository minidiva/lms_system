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

func TestGetCourseById_Success(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1, Name: "Go курс"}}, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetCourseById(context.Background(), uint(1))

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.Course.ID)
	assert.Equal(t, "Go курс", result.Course.Name)
	mockCourseRepo.AssertExpectations(t)
}

func TestGetCourseById_NotFound(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetCourseById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, entity.CourseAggregate{}, result)
}

func TestGetCourseById_DBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	dbError := errors.New("connection refused")
	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetCourseById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	assert.Equal(t, entity.CourseAggregate{}, result)
}
