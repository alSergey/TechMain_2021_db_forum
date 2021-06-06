package postgres

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
)

type ThreadRepository struct {
	conn *pgx.ConnPool
}

func NewThreadRepository(conn *pgx.ConnPool) thread.ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}

func (tr *ThreadRepository) InsertThread(thread *models.Thread) error {
	query := tr.conn.QueryRow(`
			INSERT INTO 
			thread(title, author, forum, message, slug, created)
			VALUES (
			$1, 
			COALESCE((SELECT nickname from users where nickname = $2), $2), 
			COALESCE((SELECT slug from forum where slug = $3), $3),
			$4,
			$5,
			$6)
			RETURNING id, title, author, forum, message, votes, slug, created;`,
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

func (tr *ThreadRepository) UpdateThreadBySlug(slug string, thread *models.Thread) error {
	query := tr.conn.QueryRow(`
			UPDATE thread 
			SET title=COALESCE(NULLIF($1, ''), title), 
			message=COALESCE(NULLIF($2, ''), message)
			WHERE slug=$3
			RETURNING id, title, author, forum, message, votes, slug, created;`,
		thread.Title,
		thread.Message,
		slug)

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

func (tr *ThreadRepository) UpdateThreadById(id int, thread *models.Thread) error {
	query := tr.conn.QueryRow(`
			UPDATE thread 
			SET title=COALESCE(NULLIF($1, ''), title), 
			message=COALESCE(NULLIF($2, ''), message)
			WHERE id=$3
			RETURNING id, title, author, forum, message, votes, slug, created;`,
		thread.Title,
		thread.Message,
		id)

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

func (tr *ThreadRepository) SelectThreadBySlug(slug string) (*models.Thread, error) {
	query := tr.conn.QueryRow(`
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

func (tr *ThreadRepository) SelectThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, error) {
	var query *pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1
					ORDER BY created DESC
					LIMIT $2;`,
				slug,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1
					ORDER BY created ASC
					LIMIT $2;`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1 AND created <= $2
					ORDER BY created DESC
					LIMIT $3;`,
				slug,
				params.Since,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1 AND created >= $2
					ORDER BY created ASC
					LIMIT $3;`,
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
			&thread.Created)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (tr *ThreadRepository) SelectThreadById(id int) (*models.Thread, error) {
	query := tr.conn.QueryRow(`
			SELECT id, title, author, forum, message, votes, slug, created FROM thread 
			WHERE id=$1 
			LIMIT 1;`,
		id)

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

func (tr *ThreadRepository) InsertVoteBySlug(slug string, vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			INSERT INTO 
			votes(nickname, voice, thread) 
			VALUES ($1, $2, (SELECT id from thread where slug = $3));`,
		vote.Nickname,
		vote.Voice,
		slug)
	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) UpdateVoteBySlug(slug string, vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			INSERT INTO 
			votes(nickname, voice, thread) 
			VALUES ($1, $2, (SELECT id from thread where slug = $3));`,
		vote.Nickname,
		vote.Voice,
		slug)
	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) InsertVoteById(vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			INSERT INTO 
			votes(nickname, voice, thread) 
			VALUES ($1, $2, $3);`,
		vote.Nickname,
		vote.Voice,
		vote.ThreadId)
	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) UpdateVoteById(vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			UPDATE votes 
			SET voice=$1 
			WHERE nickname=$2 and thread=$3;`,
		vote.Voice,
		vote.Nickname,
		vote.ThreadId)
	if err != nil {
		return err
	}

	return nil
}
