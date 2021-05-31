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

func (pr *PostRepository) UpdatePost(post *models.Post) error {
	query := pr.conn.QueryRow(`
			UPDATE post 
			SET message=COALESCE(NULLIF($1, ''), message),
			isEdited = CASE WHEN $1 = '' OR message = $1 THEN isEdited ELSE true END
			WHERE id = $2
			RETURNING id, parent, author, message, isEdited, forum, thread, created`,
		post.Message,
		post.Id)

	var parent sql.NullInt64
	err := query.Scan(
		&post.Id,
		&parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created)
	if err != nil {
		return err
	}

	post.Parent = null.NewIntFromNull(parent)
	return nil
}

func (pr *PostRepository) SelectPostsByFlatSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
					ORDER BY id DESC
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
					ORDER BY id ASC 
					LIMIT NULLIF($2, 0)`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND id < $2 
					ORDER BY id DESC
					LIMIT NULLIF($3, 0)`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, err
}

func (pr *PostRepository) SelectPostsByTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1)
 					ORDER BY path DESC, id DESC 
					LIMIT $2`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1) 
					ORDER BY path ASC, id ASC
					LIMIT $2`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = (SELECT id from thread where slug = $1) AND PATH < (SELECT path FROM post WHERE id = $2)
					ORDER BY path DESC, id DESC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostsByParentTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path, id`,
				slug,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path, id`,
				slug,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = (SELECT id from thread where slug = $1) AND parent IS NULL AND PATH[1] <
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path, id`,
				slug,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostsByFlatId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1
					ORDER BY id DESC
					LIMIT NULLIF($2, 0)`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1
					ORDER BY id ASC 
					LIMIT NULLIF($2, 0)`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1 AND id < $2 
					ORDER BY id DESC
					LIMIT NULLIF($3, 0)`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, err
}

func (pr *PostRepository) SelectPostsByTreeId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1 
					ORDER BY path DESC, id DESC 
					LIMIT $2`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1 
					ORDER BY path ASC, id ASC
					LIMIT $2`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE thread = $1 AND PATH < (SELECT path FROM post WHERE id = $2)
					ORDER BY path DESC, id DESC
					LIMIT $3`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostsByParentTreeId(id int, params *models.PostParams) ([]*models.Post, error) {
	var rows *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path, id`,
				id,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path, id`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] <
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path, id`,
				id,
				params.Since,
				params.Limit)
		} else {
			rows, err = pr.conn.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
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
			&post.Created)
		if err != nil {
			return nil, err
		}

		post.Parent = null.NewIntFromNull(parent)
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) SelectPostById(id int, params models.GetPostType) (*models.FullPost, error) {
	var row *pgx.Row

	switch params {
	case models.GetUser:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, u.nickname, u.fullname, u.about, u.email FROM post as p
				INNER JOIN users u on u.nickname = p.author
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetThread:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM post as p
	   		INNER JOIN thread t on t.id = p.thread
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetForum:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, f.title, f.user, f.slug, f.posts, f.threads FROM post as p
	   		INNER JOIN forum f on f.slug = p.forum
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetUserThread:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, u.nickname, u.fullname, u.about, u.email, t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM post as p
	   		INNER JOIN users u on u.nickname = p.author
				INNER JOIN thread t on t.id = p.thread
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetUserForum:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, u.nickname, u.fullname, u.about, u.email, f.title, f.user, f.slug, f.posts, f.threads FROM post as p
	   		INNER JOIN users u on u.nickname = p.author
	   		INNER JOIN forum f on f.slug = p.forum
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetThreadForum:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, f.title, f.user, f.slug, f.posts, f.threads, t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM post as p
	   		INNER JOIN forum f on f.slug = p.forum
	   		INNER JOIN thread t on t.id = p.thread
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetUserThreadForum:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created, u.nickname, u.fullname, u.about, u.email, f.title, f.user, f.slug, f.posts, f.threads, t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM post as p
	   		INNER JOIN users u on u.nickname = p.author
	   		INNER JOIN forum f on f.slug = p.forum
	   		INNER JOIN thread t on t.id = p.thread
				WHERE p.id = $1
				LIMIT 1`,
			id)

	case models.GetPost:
		row = pr.conn.QueryRow(`
				SELECT p.id, p.parent, p.author, p.message, p.isEdited, p.forum, p.thread, p.created FROM post as p
				WHERE p.id = $1
				LIMIT 1`,
			id)

	}

	var parent sql.NullInt64
	fullPost := &models.FullPost{
		Post: &models.Post{},
	}
	switch params {
	case models.GetUser:
		fullPost.Author = &models.User{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Author.NickName,
			&fullPost.Author.FullName,
			&fullPost.Author.About,
			&fullPost.Author.Email)
		if err != nil {
			return nil, err
		}

	case models.GetThread:
		fullPost.Thread = &models.Thread{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Thread.Id,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Created)
		if err != nil {
			return nil, err
		}

	case models.GetForum:
		fullPost.Forum = &models.Forum{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads)
		if err != nil {
			return nil, err
		}

	case models.GetUserThread:
		fullPost.Author = &models.User{}
		fullPost.Thread = &models.Thread{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Author.NickName,
			&fullPost.Author.FullName,
			&fullPost.Author.About,
			&fullPost.Author.Email,
			&fullPost.Thread.Id,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Created)
		if err != nil {
			return nil, err
		}

	case models.GetUserForum:
		fullPost.Author = &models.User{}
		fullPost.Forum = &models.Forum{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Author.NickName,
			&fullPost.Author.FullName,
			&fullPost.Author.About,
			&fullPost.Author.Email,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads)
		if err != nil {
			return nil, err
		}

	case models.GetThreadForum:
		fullPost.Thread = &models.Thread{}
		fullPost.Forum = &models.Forum{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
			&fullPost.Thread.Id,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Created)
		if err != nil {
			return nil, err
		}

	case models.GetUserThreadForum:
		fullPost.Author = &models.User{}
		fullPost.Thread = &models.Thread{}
		fullPost.Forum = &models.Forum{}
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created,
			&fullPost.Author.NickName,
			&fullPost.Author.FullName,
			&fullPost.Author.About,
			&fullPost.Author.Email,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
			&fullPost.Thread.Id,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Created)
		if err != nil {
			return nil, err
		}

	case models.GetPost:
		err := row.Scan(
			&fullPost.Post.Id,
			&parent,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Created)
		if err != nil {
			return nil, err
		}
	}

	fullPost.Post.Parent = null.NewIntFromNull(parent)
	return fullPost, nil
}
