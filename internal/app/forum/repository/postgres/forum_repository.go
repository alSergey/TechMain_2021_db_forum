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
