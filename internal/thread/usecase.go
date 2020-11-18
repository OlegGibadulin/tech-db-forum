package thread

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) *errors.Error
	Update(threadSlugOrID string, threadData *models.Thread) (*models.Thread, *errors.Error)
	GetBySlug(threadSlug string) (*models.Thread, *errors.Error)
	GetByID(threadID uint64) (*models.Thread, *errors.Error)
	GetBySlugOrID(threadSlugOrID string) (*models.Thread, *errors.Error)
	ListByForum(forumSlug string, filter *models.Filter) ([]*models.Thread, *errors.Error)
}
