package forum

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type ForumRepository interface {
	Insert(forum *models.Forum) error
	SelectBySlug(slug string) (*models.Forum, error)
}
