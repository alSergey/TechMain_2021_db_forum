package usecase

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ForumUsecase struct {
	forumRepo forum.ForumRepository
}

func NewForumUsecase(forumRepo forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRepo: forumRepo,
	}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, *errors.Error) {
	err := fu.forumRepo.InsertForum(forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.ForumNotExist)
			}

			if pgErr.Code == "23505" {
				existForum, err := fu.forumRepo.SelectForumBySlug(forum.Slug)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}

				return existForum, errors.Cause(errors.ForumCreateConflict)
			}
		}

		return nil, errors.UnexpectedInternal(err)
	}

	return nil, nil
}

func (fu *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, *errors.Error) {
	forum, err := fu.forumRepo.SelectForumBySlug(slug)
	if err != nil {
		return nil, errors.Cause(errors.ForumNotExist)
	}

	return forum, nil
}
