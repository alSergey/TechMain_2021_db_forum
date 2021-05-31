package service

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ServiceUsecase interface {
	GetServiceStatus() (*models.Status, *errors.Error)
	ClearService() *errors.Error
}
