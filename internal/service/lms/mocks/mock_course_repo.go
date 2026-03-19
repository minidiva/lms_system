package mocks

import (
	"context"
	"lms_system/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockCourseRepo struct {
	mock.Mock
}

func (m *MockCourseRepo) CreateCourse(ctx context.Context, entity entity.Course) (id uint, err error) {
	args := m.Called(ctx, entity)

	return args.Get(0).(uint), args.Error(1)
}

func (m *MockCourseRepo) UpdateCourseById(ctx context.Context, course entity.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepo) DeleteCourseById(ctx context.Context, courseId uint) error {
	args := m.Called(ctx, courseId)
	return args.Error(0)
}

func (m *MockCourseRepo) GetCourseById(ctx context.Context, id uint) (*entity.CourseAggregate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CourseAggregate), args.Error(1)
}

func (m *MockCourseRepo) GetAllCourses(ctx context.Context) ([]entity.Course, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Course), args.Error(1)
}
