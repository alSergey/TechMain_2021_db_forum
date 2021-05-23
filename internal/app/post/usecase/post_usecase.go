package usecase

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
)

type PostUsecase struct {
	postRepo post.PostRepository
}

func NewForumUsecase(postRepo post.PostRepository) post.PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
	}
}
