package usecase

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewForumUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uu *UserUsecase) CreateUser(user *models.User) ([]*models.User, *errors.Error) {
	err := uu.userRepo.InsertUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			users, err := uu.userRepo.SelectUserByNickNameAndEmail(user.NickName, user.Email)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

			return users, errors.Cause(errors.UserCreateConflict)
		}

		return nil, errors.UnexpectedInternal(err)
	}

	return nil, nil
}

func (uu *UserUsecase) EditUser(user *models.User) *errors.Error {
	err := uu.userRepo.UpdateUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return errors.Cause(errors.UserProfileConflict)
		}

		return errors.Cause(errors.UserNotExist)
	}

	return nil
}

func (uu *UserUsecase) GetUserByNickName(nickname string) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectUserByNickName(nickname)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	return user, nil
}
