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

func (fr *ForumRepository) Insert(forum *models.Forum) error {
	query := fr.conn.QueryRow(`
			INSERT INTO 
			forum(title, "user", slug) 
			VALUES ($1, 
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

func (fr *ForumRepository) SelectBySlug(slug string) (*models.Forum, error) {
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

func (fr *ForumRepository) InsertThread(thread *models.Thread) error {
	query := fr.conn.QueryRow(`
			INSERT INTO 
			thread(title, author, forum, message, slug, created)
			VALUES ($1, 
			COALESCE((SELECT nickname from users where nickname = $2), $2), 
			$3,
			$4,
			$5,
			$6)
			RETURNING id, title, author, forum, message, votes, slug, created`,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created)

	err := query.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	if err != nil {
		return err
	}

	return nil
}

func (fr *ForumRepository) SelectThreadBySlug(slug string) (*models.Thread, error) {
	query := fr.conn.QueryRow(`
			SELECT id, title, author, forum, message, votes, slug, created FROM thread 
			WHERE slug=$1 
			LIMIT 1;`,
		slug)

	thread := &models.Thread{}
	err := query.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	if err != nil {
		return nil, err
	}

	return thread, err
}

func (fr *ForumRepository) SelectThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, error) {
	var query *pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			query, err = fr.conn.Query(`
				SELECT id, title, author, forum, message, votes, slug, created FROM thread
				WHERE forum=$1
				ORDER BY created DESC
				LIMIT $2`,
				slug,
				params.Limit)
		} else {
			query, err = fr.conn.Query(`
				SELECT id, title, author, forum, message, votes, slug, created FROM thread
				WHERE forum=$1
				ORDER BY created ASC
				LIMIT $2`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			query, err = fr.conn.Query(`
				SELECT id, title, author, forum, message, votes, slug, created FROM thread
				WHERE forum=$1 AND created <= $2
				ORDER BY created DESC
				LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		} else {
			query, err = fr.conn.Query(`
				SELECT id, title, author, forum, message, votes, slug, created FROM thread
				WHERE forum=$1 AND created >= $2
				ORDER BY created ASC
				LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer query.Close()

	var threads []*models.Thread
	for query.Next() {
		thread := &models.Thread{}
		err := query.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
