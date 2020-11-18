package thread

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) *errors.Error
	GetBySlug(slug string) (*models.Thread, *errors.Error)
}
