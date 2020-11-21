package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostUsecase interface {
	Create(posts []*models.Post, threadID uint64) *errors.Error
	Update(postID uint64, postData *models.Post) (*models.Post, *errors.Error)
	GetByID(postID uint64) (*models.Post, *errors.Error)
	CheckAuthorsExistence(posts []*models.Post) *errors.Error
}
