package thread

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type ThreadRepository interface {
	Insert(thread *models.Thread) error
	SelectBySlug(slug string) (*models.Thread, error)
}
