package mocks

import (
	"context"
	"lms_system/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

// CreateUserCourseAccess(ctx context.Context, userAccess entity.UserCourseAccess) error
// GetByUserIdAndCourseId(ctx context.Context, userId string, courseId uint) (*entity.UserCourseAccess, error)
// GetAllByUserId(ctx context.Context, userId string) ([]entity.UserCourseAccess, error)
// UpdateAccess(ctx context.Context, access *entity.UserCourseAccess) error
type MockUserCourseAccessRepo struct {
	mock.Mock
}

func (m *MockUserCourseAccessRepo) CreateUserCourseAccess(ctx context.Context, userAccess entity.UserCourseAccess) error {
	args := m.Called(ctx, userAccess)
	return args.Error(0)
}

func (m *MockUserCourseAccessRepo) GetByUserIdAndCourseId(ctx context.Context, userId string, courseId uint) (*entity.UserCourseAccess, error) {
	args := m.Called(ctx, userId, courseId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserCourseAccess), args.Error(1)
}

func (m *MockUserCourseAccessRepo) GetAllByUserId(ctx context.Context, userId string) ([]entity.UserCourseAccess, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserCourseAccess), args.Error(1)
}

func (m *MockUserCourseAccessRepo) UpdateAccess(ctx context.Context, access *entity.UserCourseAccess) error {
	args := m.Called(ctx, access)
	return args.Error(0)
}
