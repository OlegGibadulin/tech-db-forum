package user

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByNickname(nickname string) (*models.User, error)
	SelectByEmail(email string) (*models.User, error)
	SelectByPostID(postID uint64) (*models.User, error)
	SelectExistingUsersCount(nicknames []string) (int, error)
	SelectAllByNicknameOrEmail(nickname string, email string) ([]*models.User, error)
	SelectAllByForum(forumSlug string, since string, pgnt *models.Pagination) ([]*models.User, error)
}
