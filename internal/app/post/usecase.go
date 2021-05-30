package post

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type PostUsecase interface {
	CreatePost(slug string, posts []*models.Post) ([]*models.Post, *errors.Error)

	GetPostsBySlugAndParams(slug string, params *models.PostParams) ([]*models.Post, *errors.Error)
}
