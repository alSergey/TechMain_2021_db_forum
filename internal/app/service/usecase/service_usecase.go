package usecase

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewForumUsecase(serviceRepo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepo,
	}
}

func (su *ServiceUsecase) GetServiceStatus() (*models.Status, *errors.Error) {
	status, err := su.serviceRepo.SelectService()
	if err != nil {
		//fmt.Println("GetServiceStatus error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	if status.UserCount == 0 || status.ForumCount == 0 || status.ThreadCount == 0 || status.PostCount == 0 {
		//fmt.Println("GetServiceStatus empty count")
		return nil, nil
	}

	//fmt.Println("GetServiceStatus status = ", status)
	return status, nil
}

func (su *ServiceUsecase) ClearService() *errors.Error {
	err := su.serviceRepo.TruncateService()
	if err != nil {
		//fmt.Println("ClearService error = ", err)
		return errors.UnexpectedInternal(err)
	}

	//fmt.Println("ClearService end")
	return nil
}
