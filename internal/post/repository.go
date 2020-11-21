package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostRepository interface {
	Insert(posts []*models.Post, threadID uint64) error
	SelectNotExistingParentPosts(posts []*models.Post) ([]uint64, error)
}
