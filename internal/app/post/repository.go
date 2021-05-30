package post

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type PostRepository interface {
	InsertPost(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error)
}
