package usecase

import (
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/jackc/pgx"
	"strconv"
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

func (pu *PostUsecase) GetPostsBySlugAndParams(slug string, params *models.PostParams) ([]*models.Post, *errors.Error) {
	var id int
	var posts []*models.Post

	id, err := strconv.Atoi(slug)
	if err != nil {
		switch params.Sort {
		case "flat":
			posts, err = pu.postRepo.SelectPostByFlatSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		case "tree":
			posts, err = pu.postRepo.SelectPostByTreeSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		case "parent_tree":
			posts, err = pu.postRepo.SelectPostByParentTreeSlug(slug, params)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

		default:
			posts, err = pu.postRepo.SelectPostByFlatSlug(slug, params)
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
		posts, err = pu.postRepo.SelectPostByFlatId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	case "tree":
		posts, err = pu.postRepo.SelectPostByTreeId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	case "parent_tree":
		posts, err = pu.postRepo.SelectPostByParentTreeId(id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	default:
		posts, err = pu.postRepo.SelectPostByFlatId(id, params)
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
