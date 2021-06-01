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
	//fmt.Println("CreateForum input = ", forum)
	err := fu.forumRepo.InsertForum(forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				//fmt.Println("CreateForum ForumNotExist")
				return nil, errors.Cause(errors.ForumNotExist)
			}

			if pgErr.Code == "23505" {
				existForum, err := fu.forumRepo.SelectForumBySlug(forum.Slug)
				if err != nil {
					//fmt.Println("CreateForum SelectForumBySlug error = ", err)
					return nil, errors.UnexpectedInternal(err)
				}

				//fmt.Println("CreateForum SelectForumBySlug exist forum = ", existForum)
				return existForum, errors.Cause(errors.ForumCreateConflict)
			}
		}

		//fmt.Println("CreateForum SelectForumBySlug UnexpectedInternal = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	//fmt.Println("CreateForum SelectForumBySlug end")
	return nil, nil
}

func (fu *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, *errors.Error) {
	//fmt.Println("GetForumBySlug input = ", slug)
	forum, err := fu.forumRepo.SelectForumBySlug(slug)
	if err != nil {
		//fmt.Println("GetForumBySlug SelectForumBySlug ForumNotExist")
		return nil, errors.Cause(errors.ForumNotExist)
	}

	//fmt.Println("GetForumBySlug forum = ", forum)
	return forum, nil
}

func (fu *ForumUsecase) GetForumUsersBySlugAndParams(slug string, params *models.ForumParams) ([]*models.User, *errors.Error) {
	//fmt.Println("GetForumUsersBySlugAndParams input = ", slug, params)
	users, err := fu.forumRepo.SelectForumUsersBySlugAndParams(slug, params)
	if err != nil {
		//fmt.Println("GetForumUsersBySlugAndParams SelectForumUsersBySlugAndParams err = ", err)
		return nil, errors.UnexpectedInternal(err)
	}

	if len(users) == 0 {
		//fmt.Println("GetForumUsersBySlugAndParams len = 0")
		_, err := fu.forumRepo.SelectForumBySlug(slug)
		if err != nil {
			//fmt.Println("GetForumUsersBySlugAndParams ForumNotExist")
			return nil, errors.Cause(errors.ForumNotExist)
		}

		//fmt.Println("GetForumUsersBySlugAndParams empty")
		return []*models.User{}, nil
	}

	//fmt.Println("GetForumUsersBySlugAndParams users = ", users)
	return users, nil
}
