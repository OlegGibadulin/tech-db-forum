package usecases

import (
	"database/sql"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
}

func NewThreadUsecase(repo thread.ThreadRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: repo,
	}
}

func (tu *ThreadUsecase) Create(thread *models.Thread) *errors.Error {
	if thread.Slug != "" {
		anotherThread, customErr := tu.GetBySlug(thread.Slug)
		if customErr == nil {
			customErr = errors.BuildByBody(CodeThreadAlreadyExists, anotherThread)
			return customErr
		} else if customErr.Code == CodeInternalError {
			return customErr
		}
	}

	if err := tu.threadRepo.Insert(thread); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (tu *ThreadUsecase) GetBySlug(slug string) (*models.Thread, *errors.Error) {
	thread, err := tu.threadRepo.SelectBySlug(slug)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeThreadDoesNotExist, "slug", slug)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return thread, nil
}
