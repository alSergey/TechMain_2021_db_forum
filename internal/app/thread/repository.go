package thread

import "github.com/alSergey/TechMain_2021_db_forum/internal/app/models"

type ThreadRepository interface {
	InsertThread(thread *models.Thread) error
	UpdateThreadBySlug(slug string, thread *models.Thread) error
	UpdateThreadById(id int, thread *models.Thread) error

	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadsBySlugAndParams(slug string, params *models.ThreadParams) ([]*models.Thread, error)
	SelectThreadById(id int) (*models.Thread, error)

	InsertVoteBySlug(slug string, vote *models.Vote) error
	UpdateVoteBySlug(slug string, vote *models.Vote) error
	InsertVoteById(vote *models.Vote) error
	UpdateVoteById(vote *models.Vote) error
}
