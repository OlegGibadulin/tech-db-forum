package user

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type UserUsecase interface {
	Create(user *models.User) *errors.Error
	Update(newUserData *models.User) (*models.User, *errors.Error)
	GetByNickname(nickname string) (*models.User, *errors.Error)
	GetByEmail(email string) (*models.User, *errors.Error)
	ListByNicknameOrEmail(nickname string, email string) ([]*models.User, *errors.Error)
}
