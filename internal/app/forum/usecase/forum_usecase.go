package usecase

import (
	"github.com/google/uuid"
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

func (fu *ForumUsecase) Create(forum *models.Forum) (*models.Forum, *errors.Error) {
	err := fu.forumRepo.Insert(forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.ForumCreateNotExist)
			}

			if pgErr.Code == "23505" {
				existForum, err := fu.forumRepo.SelectBySlug(forum.Slug)
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

func (fu *ForumUsecase) GetBySlug(slug string) (*models.Forum, *errors.Error) {
	forum, err := fu.forumRepo.SelectBySlug(slug)
	if err != nil {
		return nil, errors.Cause(errors.ForumDetailsNotExist)
	}

	return forum, nil
}

func (fu *ForumUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	if thread.Slug == "" {
		thread.Slug = uuid.New().String()
	}

	err := fu.forumRepo.InsertThread(thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.ForumCreateThreadNotExist)
			}

			if pgErr.Code == "23505" {
				existThread, err := fu.forumRepo.SelectThreadBySlug(thread.Slug)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}

				return existThread, errors.Cause(errors.ForumCreateThreadConflict)
			}
		}

		return nil, errors.UnexpectedInternal(err)
	}

	return nil, nil
}

func (fu *ForumUsecase) GetThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, *errors.Error) {
	threads, err := fu.forumRepo.SelectThreadsBySlugAndParams(slug, params)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if len(threads) == 0 {
		_, err := fu.forumRepo.SelectBySlug(slug)
		if err != nil {
			return nil, errors.Cause(errors.ForumThreadsNotExist)
		}

		return []*models.Thread{}, nil
	}

	return threads, nil
}
