package post

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type PostRepository interface {
	InsertPost(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error)
	UpdatePost(post *models.Post) error

	SelectPostsByFlatSlug(slug string, params *models.PostParams) ([]*models.Post, error)
	SelectPostsByTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error)
	SelectPostsByParentTreeSlug(slug string, params *models.PostParams) ([]*models.Post, error)

	SelectPostsByFlatId(id int, params *models.PostParams) ([]*models.Post, error)
	SelectPostsByTreeId(id int, params *models.PostParams) ([]*models.Post, error)
	SelectPostsByParentTreeId(id int, params *models.PostParams) ([]*models.Post, error)

	SelectPostById(id int, params models.GetPostType) (*models.FullPost, error)
}
