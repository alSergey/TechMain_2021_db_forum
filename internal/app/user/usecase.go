package user

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type UserUsecase interface {
	Create(user *models.User) ([]*models.User, *errors.Error)
	Edit(user *models.User) *errors.Error

	GetByNickName(nickname string) (*models.User, *errors.Error)
}
