package postgres

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
)

type UserRepository struct {
	conn *pgx.ConnPool
}

func NewUserRepository(conn *pgx.ConnPool) user.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (ur *UserRepository) InsertUser(user *models.User) error {
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

func (ur *UserRepository) UpdateUser(user *models.User) error {
	query := ur.conn.QueryRow(`
			UPDATE users SET
			fullname=COALESCE(NULLIF($1, ''), fullname),
			about=COALESCE(NULLIF($2, ''), about),
			email=COALESCE(NULLIF($3, ''), email)
			WHERE nickname=$4
			RETURNING nickname, fullname, about, email`,
		user.FullName,
		user.About,
		user.Email,
		user.NickName)

	err := query.Scan(
		&user.NickName,
		&user.FullName,
		&user.About,
		&user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SelectUserByNickName(nickname string) (*models.User, error) {
	query := ur.conn.QueryRow(`
			SELECT nickname, fullname, about, email FROM users 
			WHERE nickname=$1 
			LIMIT 1`,
		nickname)

	user := &models.User{}
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

func (ur *UserRepository) SelectUserByNickNameAndEmail(nickname string, email string) ([]*models.User, error) {
	query, err := ur.conn.Query(`
			SELECT nickname, fullname, about, email FROM users
			WHERE nickname=$1 or email=$2
			LIMIT 2`,
		nickname,
		email)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var users []*models.User
	for query.Next() {
		user := &models.User{}
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
