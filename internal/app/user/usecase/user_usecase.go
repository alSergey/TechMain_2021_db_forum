package usecase

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user/model"
	"github.com/jackc/pgx"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewForumUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uu *UserUsecase) Create(user *model.User) ([]*model.User, *errors.Error) {
	err := uu.userRepo.Insert(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			users, err := uu.userRepo.SelectByNickNameAndEmail(user.NickName, user.Email)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

			return users, errors.Cause(errors.UserCreateExist)
		}

		return nil, errors.UnexpectedInternal(err)
	}

	return nil, nil
}

func (uu *UserUsecase) Edit(user *model.User) *errors.Error {
	err := uu.userRepo.Update(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return errors.Cause(errors.UserProfileConflict)
		}

		return errors.Cause(errors.UserProfileNotExist)
	}

	return nil
}

func (uu *UserUsecase) GetByNickName(nickname string) (*model.User, *errors.Error) {
	user, err := uu.userRepo.SelectByNickName(nickname)
	if err != nil {
		return nil, errors.Cause(errors.UserProfileNotExist)
	}

	return user, nil
}
