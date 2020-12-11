package thread

import (
	"time"

	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) *errors.Error
	Update(threadSlugOrID string, threadData *models.Thread) (*models.Thread, *errors.Error)
	GetBySlug(threadSlug string) (*models.Thread, *errors.Error)
	GetByID(threadID uint64) (*models.Thread, *errors.Error)
	GetBySlugOrID(threadSlugOrID string) (*models.Thread, *errors.Error)
	GetByPostID(postID uint64) (*models.Thread, *errors.Error)
	CheckThreadExistence(threadSlugOrID string) (uint64, *errors.Error)
	Vote(threadSlugOrID string, vote *models.Vote) (*models.Thread, *errors.Error)
	ListByForum(forumSlug string, since time.Time, pgnt *models.Pagination) ([]*models.Thread, *errors.Error)
}
