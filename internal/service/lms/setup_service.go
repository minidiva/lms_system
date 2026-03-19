package lms

import (
	"lms_system/internal/service/lms/mocks"
	"testing"

	"github.com/sirupsen/logrus"
)

func setupService(t *testing.T) (*Service, *mocks.MockCourseRepo, *mocks.MockUserCourseAccessRepo) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)

	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	service := &Service{
		repo:   mockMainRepo,
		logger: logrus.New(),
	}

	return service, mockCourseRepo, mockAccessRepo
}
