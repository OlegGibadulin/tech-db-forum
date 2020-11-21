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

func (pu *PostUsecase) Create(posts []*models.Post, threadID uint64) *errors.Error {
	if err := pu.CheckAuthorsExistence(posts); err != nil {
		return err
	}
	if err := pu.postRepo.Insert(posts, threadID); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (pu *PostUsecase) Update(postID uint64, postData *models.Post) (*models.Post, *errors.Error) {
	post, customErr := pu.GetByID(postID)
	if customErr != nil {
		return nil, customErr
	}

	if postData.Message != "" {
		post.Message = postData.Message
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

func (pu *PostUsecase) CheckAuthorsExistence(posts []*models.Post) *errors.Error {
	postsID, err := pu.postRepo.SelectNotExistingParentPosts(posts)
	if err != nil {
		return errors.New(CodeInternalError, err)
	}
	if len(postsID) != 0 {
		postID := strconv.Itoa(int(postsID[0]))
		return errors.BuildByMsg(CodeParentPostDoesNotExist, "id", postID)
	}
	return nil
}
