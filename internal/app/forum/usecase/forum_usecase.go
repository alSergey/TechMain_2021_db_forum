package usecase

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"

type ForumUsecase struct {
	forumRepo forum.ForumRepository
}

func NewForumUsecase(forumRepo forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRepo: forumRepo,
	}
}
