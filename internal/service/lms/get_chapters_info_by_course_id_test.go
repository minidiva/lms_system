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

func TestGetChaptersInfoByCourseId_Success(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)
	mockChapterRepo.On("GetChaptersByCourseId", mock.Anything, uint(1)).
		Return([]entity.Chapter{{ID: 1, Name: "Глава 1"}, {ID: 2, Name: "Глава 2"}}, nil)
	mockLessonRepo.On("GetAllLessonsByChapterId", mock.Anything, uint(1)).
		Return([]entity.Lesson{{Name: "Урок 1"}, {Name: "Урок 2"}}, nil)
	mockLessonRepo.On("GetAllLessonsByChapterId", mock.Anything, uint(2)).
		Return([]entity.Lesson{{Name: "Урок 3"}}, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetChaptersInfoByCourseId(context.Background(), uint(1))

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, []string{"Урок 1", "Урок 2"}, result[0].LessonsName)
	assert.Equal(t, []string{"Урок 3"}, result[1].LessonsName)
	mockCourseRepo.AssertExpectations(t)
	mockChapterRepo.AssertExpectations(t)
	mockLessonRepo.AssertExpectations(t)
}

func TestGetChaptersInfoByCourseId_CourseNotFound(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetChaptersInfoByCourseId(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockChapterRepo.AssertNotCalled(t, "GetChaptersByCourseId")
}

func TestGetChaptersInfoByCourseId_CourseDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	dbError := errors.New("connection refused")
	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetChaptersInfoByCourseId(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbError)
	mockChapterRepo.AssertNotCalled(t, "GetChaptersByCourseId")
}

func TestGetChaptersInfoByCourseId_ChaptersDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)

	dbError := errors.New("connection refused")
	mockChapterRepo.On("GetChaptersByCourseId", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetChaptersInfoByCourseId(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbError)
}

func TestGetChaptersInfoByCourseId_LessonsDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)
	mockChapterRepo.On("GetChaptersByCourseId", mock.Anything, uint(1)).
		Return([]entity.Chapter{{ID: 1, Name: "Глава 1"}}, nil)

	dbError := errors.New("connection refused")
	mockLessonRepo.On("GetAllLessonsByChapterId", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	result, err := service.GetChaptersInfoByCourseId(context.Background(), uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbError)
}
