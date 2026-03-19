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

func TestCreateLesson_Success(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(&entity.Chapter{ID: 1}, nil)
	mockLessonRepo.On("CreateLesson", mock.Anything, uint(1), mock.AnythingOfType("entity.Lesson")).
		Return(uint(5), nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateLesson(context.Background(), uint(1), entity.Lesson{Name: "Урок 1"})

	assert.NoError(t, err)
	assert.Equal(t, uint(5), id)
	mockChapterRepo.AssertExpectations(t)
	mockLessonRepo.AssertExpectations(t)
}

func TestCreateLesson_ChapterNotFound(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateLesson(context.Background(), uint(1), entity.Lesson{Name: "Урок 1"})

	assert.Error(t, err)
	assert.Equal(t, uint(0), id)
	assert.Contains(t, err.Error(), "not found")
	mockLessonRepo.AssertNotCalled(t, "CreateLesson")
}

func TestCreateLesson_ChapterDBError(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	dbError := errors.New("connection refused")
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateLesson(context.Background(), uint(1), entity.Lesson{Name: "Урок 1"})

	assert.Error(t, err)
	assert.Equal(t, uint(0), id)
	assert.ErrorIs(t, err, dbError)
	mockLessonRepo.AssertNotCalled(t, "CreateLesson")
}

func TestCreateLesson_CreateDBError(t *testing.T) {
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockChapterRepo.On("GetChapterById", mock.Anything, uint(1)).
		Return(&entity.Chapter{ID: 1}, nil)

	dbError := errors.New("insert failed")
	mockLessonRepo.On("CreateLesson", mock.Anything, uint(1), mock.AnythingOfType("entity.Lesson")).
		Return(uint(0), dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateLesson(context.Background(), uint(1), entity.Lesson{Name: "Урок 1"})

	assert.Error(t, err)
	assert.Equal(t, uint(0), id)
	assert.ErrorIs(t, err, dbError)
}
