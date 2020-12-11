package thread

import (
	"time"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ThreadRepository interface {
	Insert(thread *models.Thread) error
	Update(thread *models.Thread) error
	VoteByID(threadID uint64, vote *models.Vote) error
	SelectIDByID(threadID uint64) (uint64, error)
	SelectIDBySlug(slug string) (uint64, error)
	SelectBySlug(slug string) (*models.Thread, error)
	SelectByID(threadID uint64) (*models.Thread, error)
	SelectByPostID(postID uint64) (*models.Thread, error)
	SelectAllByForum(forumSlug string, since time.Time, pgnt *models.Pagination) ([]*models.Thread, error)
}
