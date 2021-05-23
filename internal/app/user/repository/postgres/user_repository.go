package postgres

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user/model"
)

type UserRepository struct {
	conn *pgx.ConnPool
}

func NewUserRepository(conn *pgx.ConnPool) user.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (ur *UserRepository) Insert(user *model.User) error {
	_, err := ur.conn.Exec(`
			INSERT INTO 
			users(nickname, fullname, about, email) 
			VALUES ($1, $2, $3, $4)`,
		user.NickName,
		user.FullName,
		user.About,
		user.Email)

	return err
}

func (ur *UserRepository) Update(user *model.User) error {
	query := ur.conn.QueryRow(`
			UPDATE users SET
			fullname=$1,
			about=$2,
			email=$3
			WHERE nickname=$4
			RETURNING nickname`,
		user.FullName,
		user.About,
		user.Email,
		user.NickName)

	nickname := ""
	err := query.Scan(
		&nickname)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SelectByNickName(nickname string) (*model.User, error) {
	query := ur.conn.QueryRow(`
			SELECT * FROM users 
			WHERE nickname=$1 
			LIMIT 1`,
		nickname)

	user := &model.User{}
	err := query.Scan(
		&user.NickName,
		&user.FullName,
		&user.About,
		&user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) SelectByNickNameAndEmail(nickname string, email string) ([]*model.User, error) {
	query, err := ur.conn.Query(`
			SELECT * FROM users
			WHERE nickname=$1 or email=$2
			LIMIT 2`,
		nickname,
		email)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var users []*model.User
	for query.Next() {
		user := &model.User{}
		err := query.Scan(
			&user.NickName,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
