package usecases

import (
	"database/sql"
	"strconv"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/post"
)

type PostUsecase struct {
	postRepo post.PostRepository
}

func NewPostUsecase(repo post.PostRepository) post.PostUsecase {
	return &PostUsecase{
		postRepo: repo,
	}
}

func (pu *PostUsecase) Create(posts []*models.Post, thread *models.Thread) *errors.Error {
	if len(posts) == 0 {
		return nil
	}
	err := pu.postRepo.Insert(posts, thread)
	if err != nil {
		if err.Error() == OnPostInsertExceptionMsgConflict {
			return errors.BuildByMsg(CodeParentPostDoesNotExist, "id", thread.ID)
		}
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (pu *PostUsecase) Update(postID uint64, postData *models.Post) (*models.Post, *errors.Error) {
	post, customErr := pu.GetByID(postID)
	if customErr != nil {
		return nil, customErr
	}

	if postData.Message != "" && postData.Message != post.Message {
		post.Message = postData.Message
		post.IsEdited = true
	}

	if err := pu.postRepo.Update(post); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return post, nil
}

func (pu *PostUsecase) GetByID(postID uint64) (*models.Post, *errors.Error) {
	post, err := pu.postRepo.SelectByID(postID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodePostDoesNotExist, "id", strconv.Itoa(int(postID)))
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return post, nil
}

func (pu *PostUsecase) ListByThread(threadID uint64, since uint64, pgnt *models.Pagination) ([]*models.Post, *errors.Error) {
	var posts []*models.Post
	var err error

	switch pgnt.Sort {
	case models.Tree:
		posts, err = pu.postRepo.SelectAllByThreadTree(threadID, since, pgnt)
	case models.ParentTree:
		posts, err = pu.postRepo.SelectAllByThreadParentTree(threadID, since, pgnt)
	default:
		posts, err = pu.postRepo.SelectAllByThreadFlat(threadID, since, pgnt)
	}
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}

	if len(posts) == 0 {
		return []*models.Post{}, nil
	}
	return posts, nil
}
