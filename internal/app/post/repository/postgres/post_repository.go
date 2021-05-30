package postgres

import (
	"database/sql"
	"fmt"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/null"
	"strings"

	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
)

type PostRepository struct {
	conn *pgx.ConnPool
}

func NewPostRepository(conn *pgx.ConnPool) post.PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

func (pr *PostRepository) InsertPost(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error) {
	query := `INSERT INTO post(parent, author, message, thread, forum) VALUES `
	var values []interface{}

	for i, post := range posts {
		value := fmt.Sprintf(
			"(NULLIF($%d, 0), $%d, $%d, $%d, $%d),",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
		)

		query += value
		values = append(values, post.Parent, post.Author, post.Message, threadId, forumSlug)
	}

	query = strings.TrimSuffix(query, ",")
	query += ` RETURNING id, parent, author, message, isEdited, forum, thread, created`

	rows, err := pr.conn.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultPosts := make([]*models.Post, 0)
	for rows.Next() {
		currentPost := &models.Post{}
		var parent sql.NullInt64

		err := rows.Scan(
			&currentPost.Id,
			&parent,
			&currentPost.Author,
			&currentPost.Message,
			&currentPost.IsEdited,
			&currentPost.Forum,
			&currentPost.Thread,
			&currentPost.Created)
		if err != nil {
			return nil, err
		}

		currentPost.Parent = null.NewIntFromNull(parent)
		resultPosts = append(resultPosts, currentPost)
	}

	if pgErr, ok := rows.Err().(pgx.PgError); ok {
		if pgErr.Code == "12345" {
			return nil, rows.Err()
		}

		if pgErr.Code == "23503" {
			return nil, rows.Err()
		}
	}

	return resultPosts, nil
}
