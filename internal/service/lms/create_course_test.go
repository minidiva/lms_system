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
)

func TestCreateCourse_Success(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	course := entity.Course{
		Name:        "Go курс",
		Description: "Описание",
	}

	mockCourseRepo.On("CreateCourse", mock.Anything, mock.AnythingOfType("entity.Course")).
		Return(uint(1), nil)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateCourse(context.Background(), course)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
	mockCourseRepo.AssertExpectations(t)
}

func TestCreateCourse_DBError(t *testing.T) {
	mockCourseRepo := new(mocks.MockCourseRepo)
	mockMainRepo := new(mocks.MockMainRepo)
	mockMainRepo.On("Course").Return(mockCourseRepo)

	dbError := errors.New("insert failed")
	mockCourseRepo.On("CreateCourse", mock.Anything, mock.AnythingOfType("entity.Course")).
		Return(uint(0), dbError)

	service := &Service{repo: mockMainRepo, logger: logrus.New()}

	id, err := service.CreateCourse(context.Background(), entity.Course{Name: "Go курс"})

	assert.Error(t, err)
	assert.Equal(t, uint(0), id)
	assert.ErrorIs(t, err, dbError)
}
