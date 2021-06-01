package postgres

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
)

type ServiceRepository struct {
	conn *pgx.ConnPool
}

func NewServiceRepository(conn *pgx.ConnPool) service.ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}

func (sr *ServiceRepository) SelectService() (*models.Status, error) {
	row := sr.conn.QueryRow(
		`SELECT * FROM
		(SELECT COUNT(*) FROM users) as userCount,
 		(SELECT COUNT(*) FROM forum) as forumCount,
		(SELECT COUNT(*) FROM thread) as threadCount,
		(SELECT COUNT(*) FROM post) as postCount`)

	status := &models.Status{}
	err := row.Scan(
		&status.UserCount,
		&status.ForumCount,
		&status.ThreadCount,
		&status.PostCount)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (sr *ServiceRepository) TruncateService() error {
	_, err := sr.conn.Exec(`TRUNCATE users, forum, thread, post, votes, forum_users`)
	if err != nil {
		return err
	}

	return nil
}
