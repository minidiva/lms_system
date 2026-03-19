package mocks

import (
	"context"
	"lms_system/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockLessonRepo struct {
	mock.Mock
}

// CreateLesson(ctx context.Context, chapterId uint, entity entity.Lesson) (id uint, err error)
// UpdateLessonById(ctx context.Context, lesson entity.Lesson) error
// DeleteLessonById(ctx context.Context, lessonId uint) error
// GetLessonById(ctx context.Context, id uint) (*entity.Lesson, error)
// GetAllLessonsByChapterId(ctx context.Context, chapterId uint) ([]entity.Lesson, error)
func (m *MockLessonRepo) CreateLesson(ctx context.Context, chapterId uint, entity entity.Lesson) (uint, error) {
	args := m.Called(ctx, chapterId, entity)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockLessonRepo) UpdateLessonById(ctx context.Context, lesson entity.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *MockLessonRepo) DeleteLessonById(ctx context.Context, lessonId uint) error {
	args := m.Called(ctx, lessonId)
	return args.Error(0)
}

func (m *MockLessonRepo) GetLessonById(ctx context.Context, id uint) (*entity.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Lesson), args.Error(1)
}
func (m *MockLessonRepo) GetAllLessonsByChapterId(ctx context.Context, chapterId uint) ([]entity.Lesson, error) {
	args := m.Called(ctx, chapterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Lesson), args.Error(1)
}
