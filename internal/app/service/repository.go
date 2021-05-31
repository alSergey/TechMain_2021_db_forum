package service

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type ServiceRepository interface {
	SelectService() (*models.Status, error)
	TruncateService() error
}
