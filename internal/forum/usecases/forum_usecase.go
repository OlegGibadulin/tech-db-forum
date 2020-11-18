package usecases

import (
	"database/sql"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ForumUsecase struct {
	forumRepo forum.ForumRepository
}

func NewForumUsecase(repo forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRepo: repo,
	}
}

func (fu *ForumUsecase) Create(forum *models.Forum) *errors.Error {
	anotherForum, customErr := fu.GetBySlug(forum.Slug)
	if customErr == nil {
		customErr = errors.BuildByBody(CodeForumAlreadyExists, anotherForum)
		return customErr
	} else if customErr.Code == CodeInternalError {
		return customErr
	}

	if err := fu.forumRepo.Insert(forum); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (fu *ForumUsecase) GetBySlug(slug string) (*models.Forum, *errors.Error) {
	forum, err := fu.forumRepo.SelectBySlug(slug)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeForumDoesNotExist, "slug", slug)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return forum, nil
}
