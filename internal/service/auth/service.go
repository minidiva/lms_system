package auth

import (
	"lms_system/internal/domain"

	"github.com/sirupsen/logrus"
)

type Service struct {
	mainRepo       domain.MainRepositoryInterface
	logger         *logrus.Logger
	keycloakClient KeycloakClientInterface
}

func NewService(mainRepo domain.MainRepositoryInterface, logger *logrus.Logger, keycloakClient KeycloakClientInterface) *Service {
	return &Service{
		mainRepo:       mainRepo,
		logger:         logger,
		keycloakClient: keycloakClient,
	}
}
