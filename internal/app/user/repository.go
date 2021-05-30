package user

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type UserRepository interface {
	InsertUser(user *models.User) error
	UpdateUser(user *models.User) error

	SelectUserByNickName(nickname string) (*models.User, error)
	SelectUserByNickNameAndEmail(nickname string, email string) ([]*models.User, error)
}
