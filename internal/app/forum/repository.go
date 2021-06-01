package forum

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type ForumRepository interface {
	InsertForum(forum *models.Forum) error

	SelectForumBySlug(slug string) (*models.Forum, error)

	SelectForumUsersBySlugAndParams(slug string, params *models.ForumParams) ([]*models.User, error)
}
