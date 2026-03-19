package mocks

import (
	int "lms_system/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockMainRepo struct {
	mock.Mock
}

func (m *MockMainRepo) Course() int.CourseRepositoryInterface {
	args := m.Called()
	return args.Get(0).(int.CourseRepositoryInterface)
}

func (m *MockMainRepo) Chapter() int.ChapterRepositoryInterface {
	args := m.Called()
	return args.Get(0).(int.ChapterRepositoryInterface)
}

func (m *MockMainRepo) Lesson() int.LessonRepositoryInterface {
	args := m.Called()
	return args.Get(0).(int.LessonRepositoryInterface)
}

func (m *MockMainRepo) UserCourseAccess() int.UserCourseAccessInterface {
	args := m.Called()
	return args.Get(0).(int.UserCourseAccessInterface)
}

func (m *MockMainRepo) Attachment() int.AttachmentRepositoryInterface {
	args := m.Called()
	return args.Get(0).(int.AttachmentRepositoryInterface)
}
