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
	//fmt.Println("CreatePost input = ", slug, posts)
	var id int
	var thread *models.Thread

	id, err := strconv.Atoi(slug)
	if err != nil {
		//fmt.Println("CreatePost slug")
		thread, err = pu.threadRepo.SelectThreadBySlug(slug)
		if err != nil {
			//fmt.Println("CreatePost SelectThreadBySlug ThreadNotExist")
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	} else {
		//fmt.Println("CreatePost id")
		thread, err = pu.threadRepo.SelectThreadById(id)
		if err != nil {
			//fmt.Println("CreatePost SelectThreadById ThreadNotExist")
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	}

	resultPosts, err := pu.postRepo.InsertPost(thread.Id, thread.Forum, posts)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "12345" {
				//fmt.Println("CreatePost InsertPost PostWrongThread")
				return nil, errors.Cause(errors.PostWrongThread)
			}

			if pgErr.Code == "23503" {
				//fmt.Println("CreatePost InsertPost UserNotExist")
				return nil, errors.Cause(errors.UserNotExist)
			}
		}

		//fmt.Println("CreatePost UnexpectedInternal error = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	//fmt.Println("CreatePost result posts = ", resultPosts)
	return resultPosts, nil
}

func (pu *PostUsecase) UpdatePost(post *models.Post) (*models.Post, *errors.Error) {
	//fmt.Println("UpdatePost input = ", post)
	err := pu.postRepo.UpdatePost(post)
	if err != nil {
		//fmt.Println("UpdatePost UpdatePost PostNotExist")
		return nil, errors.Cause(errors.PostNotExist)
	}

	//fmt.Println("UpdatePost result post = ", post)
	return post, nil
}

func (pu *PostUsecase) GetPostsBySlugAndParams(slug string, params *models.PostParams) ([]*models.Post, *errors.Error) {
	//fmt.Println("GetPostsBySlugAndParams input = ", slug, params)
	var id int
	var posts []*models.Post

	id, err := strconv.Atoi(slug)
	if err != nil {
		//fmt.Println("GetPostsBySlugAndParams slug")
		switch params.Sort {
		case "flat":
			posts, err = pu.postRepo.SelectPostsByFlatSlug(slug, params)
			if err != nil {
				//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatSlug error = ", err)
				return nil, errors.UnexpectedInternal(err)
			}
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatSlug")

		case "tree":
			posts, err = pu.postRepo.SelectPostsByTreeSlug(slug, params)
			if err != nil {
				//fmt.Println("GetPostsBySlugAndParams SelectPostsByTreeSlug error = ", err)
				return nil, errors.UnexpectedInternal(err)
			}
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByTreeSlug")

		case "parent_tree":
			posts, err = pu.postRepo.SelectPostsByParentTreeSlug(slug, params)
			if err != nil {
				//fmt.Println("GetPostsBySlugAndParams SelectPostsByParentTreeSlug error = ", err)
				return nil, errors.UnexpectedInternal(err)
			}
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByParentTreeSlug")

		default:
			posts, err = pu.postRepo.SelectPostsByFlatSlug(slug, params)
			if err != nil {
				//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatSlug error = ", err)
				return nil, errors.UnexpectedInternal(err)
			}
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatSlug")

		}

		if len(posts) == 0 {
			//fmt.Println("GetPostsBySlugAndParams slug len 0")
			_, err := pu.threadRepo.SelectThreadBySlug(slug)
			if err != nil {
				//fmt.Println("GetPostsBySlugAndParams SelectThreadBySlug error = ", err)
				return nil, errors.Cause(errors.ThreadNotExist)
			}

			//fmt.Println("GetPostsBySlugAndParams empty")
			return []*models.Post{}, nil
		}

		//fmt.Println("GetPostsBySlugAndParams posts = ", posts)
		return posts, nil
	}

	switch params.Sort {
	case "flat":
		posts, err = pu.postRepo.SelectPostsByFlatId(id, params)
		if err != nil {
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatId error = ", err)
			return nil, errors.UnexpectedInternal(err)
		}
		//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatId")

	case "tree":
		posts, err = pu.postRepo.SelectPostsByTreeId(id, params)
		if err != nil {
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByTreeId error = ", err)
			return nil, errors.UnexpectedInternal(err)
		}
		//fmt.Println("GetPostsBySlugAndParams SelectPostsByTreeId")

	case "parent_tree":
		posts, err = pu.postRepo.SelectPostsByParentTreeId(id, params)
		if err != nil {
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByParentTreeId error = ", err)
			return nil, errors.UnexpectedInternal(err)
		}
		//fmt.Println("GetPostsBySlugAndParams SelectPostsByParentTreeId")

	default:
		posts, err = pu.postRepo.SelectPostsByFlatId(id, params)
		if err != nil {
			//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatId error = ", err)
			return nil, errors.UnexpectedInternal(err)
		}
		//fmt.Println("GetPostsBySlugAndParams SelectPostsByFlatId")

	}

	if len(posts) == 0 {
		_, err := pu.threadRepo.SelectThreadById(id)
		if err != nil {
			//fmt.Println("GetPostsBySlugAndParams SelectThreadById error = ", err)
			return nil, errors.Cause(errors.ThreadNotExist)
		}

		//fmt.Println("GetPostsBySlugAndParams empty")
		return []*models.Post{}, nil
	}

	//fmt.Println("GetPostsBySlugAndParams posts = ", posts)
	return posts, nil
}

func (pu *PostUsecase) GetPost(id int, params *models.FullPostParams) (*models.FullPost, *errors.Error) {
	//fmt.Println("GetPost input = ", id, params)
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
	//fmt.Println("GetPost param type = ", paramType)

	fullPosts, err := pu.postRepo.SelectPostById(id, paramType)
	if err != nil {
		//fmt.Println("GetPost SelectPostById error = ", err)
		return nil, errors.Cause(errors.PostNotExist)
	}

	//fmt.Println("GetPost full posts = ", fullPosts)
	return fullPosts, nil
}
