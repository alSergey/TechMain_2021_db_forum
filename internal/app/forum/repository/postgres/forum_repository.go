package postgres

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
)

type ForumRepository struct {
	conn *pgx.ConnPool
}

func NewForumRepository(conn *pgx.ConnPool) forum.ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}

func (fr *ForumRepository) InsertForum(forum *models.Forum) error {
	query := fr.conn.QueryRow(`
			INSERT INTO 
			forum(title, "user", slug) 
			VALUES (
			$1, 
			COALESCE((SELECT nickname FROM users WHERE nickname = $2), $2), 
			$3)
			RETURNING title, "user", slug, posts, threads`,
		forum.Title,
		forum.User,
		forum.Slug)

	err := query.Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads)
	if err != nil {
		return err
	}

	return nil
}

func (fr *ForumRepository) SelectForumBySlug(slug string) (*models.Forum, error) {
	query := fr.conn.QueryRow(`
			SELECT title, "user", slug, posts, threads FROM forum 
			WHERE slug=$1 
			LIMIT 1;`,
		slug)

	forum := &models.Forum{}
	err := query.Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads)
	if err != nil {
		return nil, err
	}

	return forum, err
}

func (fr *ForumRepository) SelectForumUsersBySlugAndParams(slug string, params *models.ForumParams) ([]*models.User, error) {
	var query *pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			query, err = fr.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE forum=$1
					ORDER BY nickname DESC
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		} else {
			query, err = fr.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE forum=$1
					ORDER BY nickname ASC
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			query, err = fr.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE forum=$1 AND nickname < $2
					ORDER BY nickname DESC
					LIMIT NULLIF($3, 0)`,
				slug,
				params.Since,
				params.Limit)
		} else {
			query, err = fr.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE forum=$1 AND nickname > $2
					ORDER BY nickname ASC
					LIMIT NULLIF($3, 0)`,
				slug,
				params.Since,
				params.Limit)
		}
	}

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
