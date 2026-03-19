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

func TestDeleteChapterById_Success(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(&entity.Chapter{ID: 1}, nil)
	mockChapterRepo.On("DeleteChapterById", mock.Anything, uint(1)).
		Return(nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteChapterById(context.Background(), uint(1))

	assert.NoError(t, err)
	mockChapterRepo.AssertExpectations(t)
}

func TestDeleteChapterById_NotFound(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteChapterById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockChapterRepo.AssertNotCalled(t, "DeleteChapterById")
}

func TestDeleteChapterById_GetChapterDBError(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	dbError := errors.New("connection refused")
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteChapterById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	mockChapterRepo.AssertNotCalled(t, "DeleteChapterById")
}

func TestDeleteChapterById_DeleteDBError(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(&entity.Chapter{ID: 1}, nil)

	dbError := errors.New("delete failed")
	mockChapterRepo.On("DeleteChapterById", mock.Anything, uint(1)).
		Return(dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.DeleteChapterById(context.Background(), uint(1))

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
}
