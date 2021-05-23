package user

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/user/model"

type UserRepository interface {
	Insert(user *model.User) error
	Update(user *model.User) error

	SelectByNickName(nickname string) (*model.User, error)
	SelectByNickNameAndEmail(nickname string, email string) ([]*model.User, error)
}
