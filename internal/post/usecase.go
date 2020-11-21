package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostUsecase interface {
	Create(posts []*models.Post, threadID uint64) *errors.Error
	CheckAuthorsExistence(posts []*models.Post) *errors.Error
}
