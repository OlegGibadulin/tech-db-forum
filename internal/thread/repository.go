package thread

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type ThreadRepository interface {
	Insert(thread *models.Thread) error
	Update(thread *models.Thread) error
	SelectBySlug(slug string) (*models.Thread, error)
	SelectByID(threadID uint64) (*models.Thread, error)
	SelectAllByForum(forumSlug string, filter *models.Filter) ([]*models.Thread, error)
}
