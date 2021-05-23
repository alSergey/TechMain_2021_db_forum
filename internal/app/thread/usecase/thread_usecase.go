package usecase

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
}

func NewForumUsecase(threadRepo thread.ThreadRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
	}
}
