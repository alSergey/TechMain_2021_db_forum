package usecase

import (
	"strconv"

	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type PostUsecase struct {
	postRepo   post.PostRepository
	threadRepo thread.ThreadRepository
}

func NewForumUsecase(postRepo post.PostRepository, threadRepo thread.ThreadRepository) post.PostUsecase {
	return &PostUsecase{
		postRepo:   postRepo,
		threadRepo: threadRepo,
	}
}

func (pu *PostUsecase) CreatePost(slug string, posts []*models.Post) ([]*models.Post, *errors.Error) {
	var id int
	var thread *models.Thread

	id, err := strconv.Atoi(slug)
	if err != nil {
		thread, err = pu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	} else {
		thread, err = pu.threadRepo.SelectThreadById(id)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	}

	resultPosts, err := pu.postRepo.InsertPost(thread.Id, thread.Forum, posts)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "12345" {
				return nil, errors.Cause(errors.PostWrongThread)
			}

			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.UserNotExist)
			}
		}

		return nil, errors.UnexpectedInternal(err)
	}

	return resultPosts, nil
}

func (pu *PostUsecase) UpdatePost(post *models.Post) (*models.Post, *errors.Error) {
	err := pu.postRepo.UpdatePost(post)
	if err != nil {
		return nil, errors.Cause(errors.PostNotExist)
	}

	return post, nil
}

func (pu *PostUsecase) GetPostsBySlugAndParams(slug string, params *models.PostParams) ([]*models.Post, *errors.Error) {
	var id int
	var posts []*models.Post

	id, err := strconv.Atoi(slug)
	if err != nil {
		switch params.Sort {
		case "flat":
			posts, err = pu.postRepo.SelectPostsByFlatSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		case "tree":
			posts, err = pu.postRepo.SelectPostsByTreeSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		case "parent_tree":
			posts, err = pu.postRepo.SelectPostsByParentTreeSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		default:
			posts, err = pu.postRepo.SelectPostsByFlatSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		}

		if len(posts) == 0 {
			_, err := pu.threadRepo.SelectThreadBySlug(slug)
			if err != nil {
				return nil, errors.Cause(errors.ThreadNotExist)
			}

			return []*models.Post{}, nil
		}

		return posts, nil
	}

	switch params.Sort {
	case "flat":
		posts, err = pu.postRepo.SelectPostsByFlatId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	case "tree":
		posts, err = pu.postRepo.SelectPostsByTreeId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	case "parent_tree":
		posts, err = pu.postRepo.SelectPostsByParentTreeId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	default:
		posts, err = pu.postRepo.SelectPostsByFlatId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	}

	if len(posts) == 0 {
		_, err := pu.threadRepo.SelectThreadById(id)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		return []*models.Post{}, nil
	}

	return posts, nil
}

func (pu *PostUsecase) GetPost(id int, params *models.FullPostParams) (*models.FullPost, *errors.Error) {
	paramType := models.GetPost

	if params.Forum {
		paramType = models.GetForum
	}

	if params.Thread {
		paramType = models.GetThread

		if params.Forum {
			paramType = models.GetThreadForum
		}
	}

	if params.User {
		paramType = models.GetUser

		if params.Thread {
			paramType = models.GetUserThread
		}

		if params.Forum {
			paramType = models.GetUserForum
		}

		if params.Thread && params.Forum {
			paramType = models.GetUserThreadForum
		}
	}

	fullPosts, err := pu.postRepo.SelectPostById(id, paramType)
	if err != nil {
		return nil, errors.Cause(errors.PostNotExist)
	}

	return fullPosts, nil
}
