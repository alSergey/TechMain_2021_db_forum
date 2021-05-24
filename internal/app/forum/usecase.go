package forum

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ForumUsecase interface {
	Create(forum *models.Forum) (*models.Forum, *errors.Error)

	GetBySlug(slug string) (*models.Forum, *errors.Error)


	CreateThread(thread *models.Thread) (*models.Thread, *errors.Error)

	GetThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, *errors.Error)
}
