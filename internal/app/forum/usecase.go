package forum

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, *errors.Error)

	GetForumBySlug(slug string) (*models.Forum, *errors.Error)

	GetForumUsersBySlugAndParams(slug string, params *models.ForumParams) ([]*models.User, *errors.Error)
}
