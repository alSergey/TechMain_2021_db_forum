package forum

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type ForumRepository interface {
	Insert(forum *models.Forum) error

	SelectBySlug(slug string) (*models.Forum, error)

	InsertThread(thread *models.Thread) error

	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, error)
}
