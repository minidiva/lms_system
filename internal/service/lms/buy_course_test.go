package lms

import (
	"context"
	"errors"
	"lms_system/internal/domain/dto"
	"lms_system/internal/domain/entity"
	"lms_system/internal/service/lms/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestBuyCourse_Success(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)

	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{
			Course: entity.Course{ID: 1},
		}, nil)

	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	mockAccessRepo.On("CreateUserCourseAccess", mock.Anything, mock.AnythingOfType("entity.UserCourseAccess")).
		Return(nil)

	logger := logrus.New()
	service := &Service{
		repo:   mockMainRepo,
		logger: logger,
	}

	request := dto.BuyCourseRequest{
		UserUUID: "user-1",
		CourseId: 1,
	}

	err := service.BuyCourse(context.Background(), request)

	assert.NoError(t, err)
	mockMainRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAccessRepo.AssertExpectations(t)
}

func TestBuyCourse_NotFound(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)

	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	logger := logrus.New()
	service := &Service{
		repo:   mockMainRepo,
		logger: logger,
	}

	err := service.BuyCourse(context.Background(), dto.BuyCourseRequest{
		UserUUID: "user-1",
		CourseId: 1,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Проверяем что дальше курса метод не пошёл
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
	mockAccessRepo.AssertNotCalled(t, "CreateUserCourseAccess")
}

// Ошибка БД при получении курса
func TestBuyCourse_CourseDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	dbError := errors.New("connection refused")
	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(nil, dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.BuyCourse(context.Background(), dto.BuyCourseRequest{
		UserUUID: "user-1",
		CourseId: 1,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
	mockAccessRepo.AssertNotCalled(t, "GetByUserIdAndCourseId")
	mockAccessRepo.AssertNotCalled(t, "CreateUserCourseAccess")
}

// Пользователь уже имеет доступ
func TestBuyCourse_AlreadyHasAccess(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(1)).
		Return(&entity.UserCourseAccess{UserID: "user-1", CourseID: 1}, nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.BuyCourse(context.Background(), dto.BuyCourseRequest{
		UserUUID: "user-1",
		CourseId: 1,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already has access")
	mockAccessRepo.AssertNotCalled(t, "CreateUserCourseAccess")
}

// Ошибка БД при создании доступа
func TestBuyCourse_CreateAccessDBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockAccessRepo := new(mocks.MockUserCourseAccessRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)
	mockMainRepo.On("UserCourseAccess").Return(mockAccessRepo)

	mockCourseRepo.On("GetCourseById", mock.Anything, uint(1)).
		Return(&entity.CourseAggregate{Course: entity.Course{ID: 1}}, nil)
	mockAccessRepo.On("GetByUserIdAndCourseId", mock.Anything, "user-1", uint(1)).
		Return(nil, gorm.ErrRecordNotFound)

	dbError := errors.New("insert failed")
	mockAccessRepo.On("CreateUserCourseAccess", mock.Anything, mock.AnythingOfType("entity.UserCourseAccess")).
		Return(dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	err := service.BuyCourse(context.Background(), dto.BuyCourseRequest{
		UserUUID: "user-1",
		CourseId: 1,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbError)
}
