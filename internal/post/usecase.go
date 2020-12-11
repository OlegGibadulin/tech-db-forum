package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostUsecase interface {
	Create(posts []*models.Post, thread *models.Thread) *errors.Error
	Update(postID uint64, postData *models.Post) (*models.Post, *errors.Error)
	GetByID(postID uint64) (*models.Post, *errors.Error)
	CheckAuthorsExistence(posts []*models.Post) *errors.Error
	ListByThread(threadID uint64, since uint64, pgnt *models.Pagination) ([]*models.Post, *errors.Error)
}
