package user

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type UserRepository interface {
	Insert(user *models.User) error
	Update(user *models.User) error

	SelectByNickName(nickname string) (*models.User, error)
	SelectByNickNameAndEmail(nickname string, email string) ([]*models.User, error)
}
