package thread

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, *errors.Error)
	UpdateThread(thread *models.Thread) (*models.Thread, *errors.Error)

	GetThread(slug string) (*models.Thread, *errors.Error)
	GetThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, *errors.Error)

	Vote(slug string, vote *models.Vote) (*models.Thread, *errors.Error)
}
