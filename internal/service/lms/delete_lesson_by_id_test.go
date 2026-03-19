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

func TestDeleteLessonById_Success(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1}, nil)
	mockLessonRepo.On("DeleteLessonById", mock.Anything, uint(1)).
		Return(nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteLessonById(context.Background(), uint(1))

	assert.NoError(t, err)
	mockLessonRepo.AssertExpectations(t)
}

func TestDeleteLessonById_NotFound(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteLessonById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockLessonRepo.AssertNotCalled(t, "DeleteLessonById")
}

func TestDeleteLessonById_GetLessonDBError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	dbError := errors.New("connection refused")
	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteLessonById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	mockLessonRepo.AssertNotCalled(t, "DeleteLessonById")
}

func TestDeleteLessonById_DeleteDBError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1}, nil)

	dbError := errors.New("delete failed")
	mockLessonRepo.On("DeleteLessonById", mock.Anything, uint(1)).
		Return(dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteLessonById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
}
