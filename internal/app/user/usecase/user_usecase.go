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
	//fmt.Println("CreateUser input = ", user)

	err := uu.userRepo.InsertUser(user)
	if err != nil {
		//fmt.Println("CreateUser insertUser error = ", err)

		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			users, err := uu.userRepo.SelectUserByNickNameAndEmail(user.NickName, user.Email)
			if err != nil {
				//fmt.Println("CreateUser SelectUserByNickNameAndEmail error = ", err)
				return nil, errors.UnexpectedInternal(err)
			}

			//fmt.Println("CreateUser UserCreateConflict = ", users)
			return users, errors.Cause(errors.UserCreateConflict)
		}

		//fmt.Println("CreateUser UnexpectedInternal error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	//fmt.Println("CreateUser end")
	return nil, nil
}

func (uu *UserUsecase) EditUser(user *models.User) *errors.Error {
	//fmt.Println("EditUser input = ", user)

	err := uu.userRepo.UpdateUser(user)
	if err != nil {
		//fmt.Println("EditUser UpdateUser error = ", err)

		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			//fmt.Println("EditUser UserProfileConflict")
			return errors.Cause(errors.UserProfileConflict)
		}

		//fmt.Println("EditUser UserNotExist")
		return errors.Cause(errors.UserNotExist)
	}

	//fmt.Println("EditUser end")
	return nil
}

func (uu *UserUsecase) GetUserByNickName(nickname string) (*models.User, *errors.Error) {
	//fmt.Println("GetUserByNickName input = ", nickname)

	user, err := uu.userRepo.SelectUserByNickName(nickname)
	if err != nil {
		//fmt.Println("GetUserByNickName SelectUserByNickName error = ", err)
		//fmt.Println("GetUserByNickName UserNotExist")
		return nil, errors.Cause(errors.UserNotExist)
	}

	//fmt.Println("GetUserByNickName end")
	return user, nil
}
