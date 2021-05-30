package post

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type PostRepository interface {
	InsertPost(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error)

	SelectPostByFlatSlug(slug string, params *models.PostParams) ([]*models.Post, error)
	SelectPostByTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error)
	SelectPostByParentTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error)

	SelectPostByFlatId(id int, params *models.PostParams) ([]*models.Post, error)
	SelectPostByTreeId(id int, params *models.PostParams) ([]*models.Post, error)
	SelectPostByParentTreeId(id int, params *models.PostParams) ([]*models.Post, error)
}
