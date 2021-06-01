package usecase

import (
	"strconv"

	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/uuid"
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
	//fmt.Println("CreateThread input = ", thread)
	if thread.Slug == "" {
		//fmt.Println("CreateThread create slug")
		thread.Slug = uuid.CreateSlug()
	}

	err := tu.threadRepo.InsertThread(thread)
	if err != nil {
		//fmt.Println("CreateThread InsertThread error = ", err)
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				//fmt.Println("CreateThread ThreadNotExist")
				return nil, errors.Cause(errors.ThreadNotExist)
			}

			if pgErr.Code == "23505" {
				existThread, err := tu.threadRepo.SelectThreadBySlug(thread.Slug)
				if err != nil {
					//fmt.Println("CreateThread SelectThreadBySlug error = ", err)
					return nil, errors.UnexpectedInternal(err)
				}

				//fmt.Println("CreateThread ForumCreateThreadConflict")
				return existThread, errors.Cause(errors.ForumCreateThreadConflict)
			}
		}

		//fmt.Println("CreateThread UnexpectedInternal error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	//fmt.Println("CreateThread end")
	return nil, nil
}

func (tu *ThreadUsecase) UpdateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	//fmt.Println("UpdateThread input = ", thread)

	id, err := strconv.Atoi(thread.Forum)
	if err != nil {
		//fmt.Println("UpdateThread slug")
		err = tu.threadRepo.UpdateThreadBySlug(thread.Forum, thread)
		if err != nil {
			//fmt.Println("UpdateThread ThreadNotExist")
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		//fmt.Println("UpdateThread thread = ", thread)
		return thread, nil
	}

	//fmt.Println("UpdateThread id")
	err = tu.threadRepo.UpdateThreadById(id, thread)
	if err != nil {
		//fmt.Println("UpdateThread ThreadNotExist")
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	//fmt.Println("UpdateThread thread = ", thread)
	return thread, nil
}

func (tu *ThreadUsecase) GetThread(slug string) (*models.Thread, *errors.Error) {
	//fmt.Println("GetThread input = ", slug)
	id, err := strconv.Atoi(slug)
	if err != nil {
		//fmt.Println("GetThread slug")
		thread, err := tu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			//fmt.Println("GetThread ThreadNotExist")
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		//fmt.Println("GetThread thread = ", thread)
		return thread, nil
	}

	//fmt.Println("GetThread id")
	thread, err := tu.threadRepo.SelectThreadById(id)
	if err != nil {
		//fmt.Println("GetThread ThreadNotExist")
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	//fmt.Println("GetThread thread = ", thread)
	return thread, nil
}

func (tu *ThreadUsecase) GetThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, *errors.Error) {
	//fmt.Println("GetThreadsBySlugAndParams input = ", slug, params)
	threads, err := tu.threadRepo.SelectThreadsBySlugAndParams(slug, params)
	if err != nil {
		//fmt.Println("GetThreadsBySlugAndParams SelectThreadsBySlugAndParams error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	if len(threads) == 0 {
		_, err := tu.forumRepo.SelectForumBySlug(slug)
		if err != nil {
			//fmt.Println("GetThreadsBySlugAndParams SelectForumBySlug ForumNotExist")
			return nil, errors.Cause(errors.ForumNotExist)
		}

		//fmt.Println("GetThreadsBySlugAndParams SelectForumBySlug empty")
		return []*models.Thread{}, nil
	}

	//fmt.Println("GetThreadsBySlugAndParams threads = ", threads)
	return threads, nil
}

func (tu *ThreadUsecase) Vote(slug string, vote *models.Vote) (*models.Thread, *errors.Error) {
	//fmt.Println("Vote input = ", slug, vote)
	id, err := strconv.Atoi(slug)
	if err != nil {
		//fmt.Println("Vote slug")
		err = tu.threadRepo.InsertVoteBySlug(slug, vote)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok {
				if pgErr.Code == "23505" {
					_ = tu.threadRepo.UpdateVoteBySlug(slug, vote)

					thread, err := tu.threadRepo.SelectThreadBySlug(slug)
					if err != nil {
						//fmt.Println("Vote SelectThreadBySlug error = ", err)
						return nil, errors.UnexpectedInternal(err)
					}

					//fmt.Println("Vote SelectThreadBySlug thread = ", thread)
					return thread, nil
				}

				if pgErr.Code == "23503" {
					//fmt.Println("Vote UserNotExist")
					return nil, errors.Cause(errors.UserNotExist)
				}

				if pgErr.Code == "23502" {
					//fmt.Println("Vote ThreadNotExist")
					return nil, errors.Cause(errors.ThreadNotExist)
				}
			}

			//fmt.Println("Vote UnexpectedInternal = ", err)
			return nil, errors.UnexpectedInternal(err)
		}

		thread, err := tu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			//fmt.Println("Vote SelectThreadBySlug error = ", err)
			return nil, errors.UnexpectedInternal(err)
		}

		//fmt.Println("Vote SelectThreadBySlug thread = ", thread)
		return thread, nil
	}

	//fmt.Println("Vote id")
	vote.ThreadId = id
	err = tu.threadRepo.InsertVoteById(vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23505" {
				_ = tu.threadRepo.UpdateVoteById(vote)

				thread, err := tu.threadRepo.SelectThreadById(id)
				if err != nil {
					//fmt.Println("Vote SelectThreadById error = ", err)
					return nil, errors.UnexpectedInternal(err)
				}

				//fmt.Println("Vote SelectThreadById thread = ", thread)
				return thread, nil
			}

			if pgErr.Code == "23503" {
				//fmt.Println("Vote UserNotExist")
				return nil, errors.Cause(errors.UserNotExist)
			}

			if pgErr.Code == "23502" {
				//fmt.Println("Vote ThreadNotExist")
				return nil, errors.Cause(errors.ThreadNotExist)
			}
		}

		//fmt.Println("Vote UnexpectedInternal = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	thread, err := tu.threadRepo.SelectThreadById(id)
	if err != nil {
		//fmt.Println("Vote SelectThreadById error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	//fmt.Println("Vote SelectThreadById thread = ", thread)
	return thread, nil
}
