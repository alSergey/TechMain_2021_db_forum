package user

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type UserUsecase interface {
	CreateUser(user *models.User) ([]*models.User, *errors.Error)
	EditUser(user *models.User) *errors.Error

	GetUserByNickName(nickname string) (*models.User, *errors.Error)
}
