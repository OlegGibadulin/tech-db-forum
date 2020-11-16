package user

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByNickname(nickname string) (*models.User, error)
	SelectByEmail(email string) (*models.User, error)
	SelectAllByNicknameOrEmail(nickname string, email string) ([]*models.User, error)
}
