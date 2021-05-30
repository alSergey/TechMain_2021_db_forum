package usecase

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strconv"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
	forumRepo  forum.ForumRepository
}

func NewForumUsecase(threadRepo thread.ThreadRepository, forumRepo forum.ForumRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}

func (tu *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	if thread.Slug == "" {
		thread.Slug = uuid.New().String()
	}

	err := tu.threadRepo.InsertThread(thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.ThreadNotExist)
			}

			if pgErr.Code == "23505" {
				existThread, err := tu.threadRepo.SelectThreadBySlug(thread.Slug)
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

func (tu *ThreadUsecase) UpdateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	id, err := strconv.Atoi(thread.Forum)
	if err != nil {
		err = tu.threadRepo.UpdateThreadBySlug(thread.Forum, thread)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		return thread, nil
	}

	err = tu.threadRepo.UpdateThreadById(id, thread)
	if err != nil {
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	return thread, nil
}

func (tu *ThreadUsecase) GetThread(slug string) (*models.Thread, *errors.Error) {
	id, err := strconv.Atoi(slug)
	if err != nil {
		thread, err := tu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		return thread, nil
	}

	thread, err := tu.threadRepo.SelectThreadById(id)
	if err != nil {
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	return thread, nil
}

func (tu *ThreadUsecase) GetThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, *errors.Error) {
	threads, err := tu.threadRepo.SelectThreadsBySlugAndParams(slug, params)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if len(threads) == 0 {
		_, err := tu.forumRepo.SelectForumBySlug(slug)
		if err != nil {
			return nil, errors.Cause(errors.ForumNotExist)
		}

		return []*models.Thread{}, nil
	}

	return threads, nil
}

func (tu *ThreadUsecase) Vote(slug string, vote *models.Vote) (*models.Thread, *errors.Error) {
	id, err := strconv.Atoi(slug)
	if err != nil {
		err = tu.threadRepo.InsertVoteBySlug(slug, vote)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok {
				if pgErr.Code == "23505" {
					_ = tu.threadRepo.UpdateVoteBySlug(slug, vote)

					thread, err := tu.threadRepo.SelectThreadBySlug(slug)
					if err != nil {
						return nil, errors.UnexpectedInternal(err)
					}

					return thread, nil
				}

				if pgErr.Code == "23503" {
					return nil, errors.Cause(errors.UserNotExist)
				}

				if pgErr.Code == "23502" {
					return nil, errors.Cause(errors.ThreadNotExist)
				}
			}

			return nil, errors.UnexpectedInternal(err)
		}

		thread, err := tu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

		return thread, nil
	}

	vote.ThreadId = id
	err = tu.threadRepo.InsertVoteById(vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23505" {
				_ = tu.threadRepo.UpdateVoteById(vote)

				thread, err := tu.threadRepo.SelectThreadById(id)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}

				return thread, nil
			}

			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.UserNotExist)
			}

			if pgErr.Code == "23502" {
				return nil, errors.Cause(errors.ThreadNotExist)
			}
		}

		return nil, errors.UnexpectedInternal(err)
	}

	thread, err := tu.threadRepo.SelectThreadById(id)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return thread, nil
}
