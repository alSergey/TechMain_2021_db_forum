package user

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user/model"
)

type UserUsecase interface {
	Create(user *model.User) ([]*model.User, *errors.Error)
	Edit(user *model.User) *errors.Error

	GetByNickName(nickname string) (*model.User, *errors.Error)
}
