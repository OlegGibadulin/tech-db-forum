package user

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type UserUsecase interface {
	Create(user *models.User) *errors.Error
	Update(nickname string, newUserData *models.User) (*models.User, *errors.Error)
	GetByNickname(nickname string) (*models.User, *errors.Error)
	GetByEmail(email string) (*models.User, *errors.Error)
	GetByPostID(postID uint64) (*models.User, *errors.Error)
	ListByNicknameOrEmail(nickname string, email string) ([]*models.User, *errors.Error)
	ListByForum(forumSlug string, since string, pgnt *models.Pagination) ([]*models.User, *errors.Error)
}
