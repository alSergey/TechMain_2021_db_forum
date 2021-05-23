package usecase

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewForumUsecase(serviceRepo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepo,
	}
}
