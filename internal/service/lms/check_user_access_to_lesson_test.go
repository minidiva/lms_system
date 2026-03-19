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

func TestCheckUserAccessToLesson_Success(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)

	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	// Шаг 1 — урок найден, привязан к главе 2
	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1, ChapterID: 2}, nil)

	// Шаг 2 — глава найдена, привязана к курсу 3
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(2)).
		Return(&entity.Chapter{ID: 2, CourseID: 3}, nil)

	// Шаг 3 — доступ есть
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(3)).
		Return(&entity.UserCourseAccess{Unlocked: true}, nil)

	service := &Service{
		repo:        mockMainRepo,
		logger:      logrus.New(),
		fileService: nil,
	}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.NoError(t, err)
	assert.True(t, hasAccess)
	mockLessonRepo.AssertExpectations(t)
	mockChapterRepo.AssertExpectations(t)
	mockAccessRepo.AssertExpectations(t)
}

func TestCheckUserAccessToLesson_NotFound(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)

	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	logger := logrus.New()
	service := &Service{
		logger: logger,
		repo:   mockMainRepo,
	}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.False(t, hasAccess)
	mockChapterRepo.AssertNotCalled(t, "GetChapterById")
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
}

// Ошибка БД при получении урока
func TestCheckUserAccessToLesson_LessonDBError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	dbError := errors.New("connection refused")
	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.Error(t, err)
	assert.False(t, hasAccess)
	assert.ErrorIs(t, err, dbError)
	mockChapterRepo.AssertNotCalled(t, "GetChapterById")
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
}

// Ошибка БД при получении главы
func TestCheckUserAccessToLesson_ChapterDBError(t *testing.T) {
	mockLessonRepo := new(mocks.MockLessonRepo)
	mockChapterRepo := new(mocks.MockChapterRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Lesson").Return(mockLessonRepo)
	mockMainRepo.On("Chapter").Return(mockChapterRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockLessonRepo.On("GetLessonById", mock.Anything, uint(1)).
		Return(&entity.Lesson{ID: 1, ChapterID: 2}, nil)

	dbError := errors.New("connection refused")
	mockChapterRepo.On("GetChapterById", mock.Anything, uint(2)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.Error(t, err)
	assert.False(t, hasAccess)
	assert.ErrorIs(t, err, dbError)
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
}

// Глава не найдена (nil, nil)
func TestCheckUserAccessToLesson_ChapterNotFound(t *testing.T) {
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
		Return(nil, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chapter not found")
	assert.False(t, hasAccess)
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
}

// Доступа нет (ErrRecordNotFound) — не ошибка, просто false
func TestCheckUserAccessToLesson_NoAccess(t *testing.T) {
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

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.NoError(t, err) // не ошибка — просто нет доступа
	assert.False(t, hasAccess)
}

// Ошибка БД при проверке доступа
func TestCheckUserAccessToLesson_AccessDBError(t *testing.T) {
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

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.Error(t, err)
	assert.False(t, hasAccess)
	assert.ErrorIs(t, err, dbError)
}

// Доступ есть но Unlocked = false
func TestCheckUserAccessToLesson_AccessLocked(t *testing.T) {
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
		Return(&entity.UserCourseAccess{Unlocked: false}, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	hasAccess, err := service.CheckUserAccessToLesson(context.Background(), "user-1", uint(1))

	assert.NoError(t, err)
	assert.False(t, hasAccess)
}
