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

func TestGetLessonById_Success(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1, Name: "Урок 1", ChapterID: 2}, nil)

	// CheckUserAccessToLesson внутри тоже вызывает GetLessonById
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(2)).
		Return(&entity.Chapter{ID: 2, CourseID: 3}, nil)
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(3)).
		Return(&entity.UserCourseAccess{Unlocked: true}, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetLessonById(context.Background(), uint(1), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Урок 1", result.Name)
}

func TestGetLessonById_NotFound(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetLessonById(context.Background(), uint(1), "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, entity.Lesson{}, result)
}

func TestGetLessonById_DBError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	dbError := errors.New("connection refused")
	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetLessonById(context.Background(), uint(1), "user-1")

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	assert.Equal(t, entity.Lesson{}, result)
}

func TestGetLessonById_AccessDenied(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1, ChapterID: 2}, nil)
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(2)).
		Return(&entity.Chapter{ID: 2, CourseID: 3}, nil)
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(3)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetLessonById(context.Background(), uint(1), "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
	assert.Equal(t, entity.Lesson{}, result)
}

func TestGetLessonById_AccessCheckError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1, ChapterID: 2}, nil)
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(2)).
		Return(&entity.Chapter{ID: 2, CourseID: 3}, nil)

	dbError := errors.New("connection refused")
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(3)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetLessonById(context.Background(), uint(1), "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check access")
	assert.Equal(t, entity.Lesson{}, result)
}
