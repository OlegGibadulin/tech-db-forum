package post

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type PostRepository interface {
	Insert(posts []*models.Post, thread *models.Thread) error
	Update(post *models.Post) error
	SelectByID(postID uint64) (*models.Post, error)
	SelectNotExistingParentPosts(posts []*models.Post) ([]uint64, error)
	SelectAllByThreadFlat(threadID uint64, since uint64, pgnt *models.Pagination) ([]*models.Post, error)
	SelectAllByThreadTree(threadID uint64, since uint64, pgnt *models.Pagination) ([]*models.Post, error)
	SelectAllByThreadParentTree(threadID uint64, since uint64, pgnt *models.Pagination) ([]*models.Post, error)
}
