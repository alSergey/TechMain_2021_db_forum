package postgres

import (
	"database/sql"
	"fmt"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/null"
	"github.com/jackc/pgx"
	"strings"
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

func (pr *PostRepository) SelectPostByFlatSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
					ORDER BY id DESC
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
					ORDER BY id ASC 
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND id < $2 
					ORDER BY id DESC
					LIMIT NULLIF($3, 0)`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND id > $2
					ORDER BY id ASC
					LIMIT NULLIF($3, 0)`,
				slug,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, err
}

func (pr *PostRepository) SelectPostByTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
 					ORDER BY path DESC, id DESC 
					LIMIT $2`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1) 
					ORDER BY path ASC, id ASC
					LIMIT $2`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND PATH < (SELECT path FROM post WHERE id = $2)
					ORDER BY path DESC, id DESC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND PATH > (SELECT path FROM post WHERE id = $2)
					ORDER BY path ASC, id ASC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostByParentTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path, id`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path, id`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL AND PATH[1] <
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path, id`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL AND PATH[1] >
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path, id`,
				slug,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostByFlatId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1
					ORDER BY id DESC
					LIMIT NULLIF($2, 0)`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1
					ORDER BY id ASC 
					LIMIT NULLIF($2, 0)`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 AND id < $2 
					ORDER BY id DESC
					LIMIT NULLIF($3, 0)`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 AND id > $2
					ORDER BY id ASC
					LIMIT NULLIF($3, 0)`,
				id,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, err
}

func (pr *PostRepository) SelectPostByTreeId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 
					ORDER BY path DESC, id DESC 
					LIMIT $2`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 
					ORDER BY path ASC, id ASC
					LIMIT $2`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 AND PATH < (SELECT path FROM post WHERE id = $2)
					ORDER BY path DESC, id DESC
					LIMIT $3`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE thread = $1 AND PATH > (SELECT path FROM post WHERE id = $2)
					ORDER BY path ASC, id ASC
					LIMIT $3`,
				id,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostByParentTreeId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path, id`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path, id`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] <
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path, id`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created, path FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] >
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path, id`,
				id,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var parent sql.NullInt64

		err = rows.Scan(
			&post.Id,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.Path)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}
