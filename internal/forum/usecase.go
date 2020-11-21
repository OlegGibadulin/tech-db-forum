package forum

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ForumUsecase interface {
	Create(forum *models.Forum) *errors.Error
	GetBySlug(slug string) (*models.Forum, *errors.Error)
	GetByPostID(postID uint64) (*models.Forum, *errors.Error)
}
