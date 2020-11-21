package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostRepository interface {
	Insert(posts []*models.Post, threadID uint64) error
	Update(post *models.Post) error
	SelectByID(postID uint64) (*models.Post, error)
	SelectNotExistingParentPosts(posts []*models.Post) ([]uint64, error)
}
