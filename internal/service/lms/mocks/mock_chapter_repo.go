package mocks

import (
	"context"
	"lms_system/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

// CreateChapter(ctx context.Context, courseId uint, entity *entity.Chapter) (id uint, err error)
// UpdateChapterById(ctx context.Context, chapter *entity.Chapter) error
// DeleteChapterById(ctx context.Context, chapterId uint) error
// GetChapterById(ctx context.Context, id uint) (*entity.Chapter, error)
// GetChaptersByCourseId(ctx context.Context, id uint) ([]entity.Chapter, error)

type MockChapterRepo struct {
	mock.Mock
}

func (m *MockChapterRepo) CreateChapter(ctx context.Context, courseId uint, entity *entity.Chapter) (uint, error) {
	args := m.Called(ctx, courseId, entity)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockChapterRepo) UpdateChapterById(ctx context.Context, chapter *entity.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterRepo) DeleteChapterById(ctx context.Context, chapterId uint) error {
	args := m.Called(ctx, chapterId)
	return args.Error(0)
}

func (m *MockChapterRepo) GetChapterById(ctx context.Context, id uint) (*entity.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Chapter), args.Error(1)
}

func (m *MockChapterRepo) GetChaptersByCourseId(ctx context.Context, id uint) ([]entity.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Chapter), args.Error(1)
}
