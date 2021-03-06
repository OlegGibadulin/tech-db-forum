package usecases

import (
	"database/sql"
	"time"

	"strconv"

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

func (tu *ThreadUsecase) Update(threadSlugOrID string, threadData *models.Thread) (*models.Thread, *errors.Error) {
	thread, customErr := tu.GetBySlugOrID(threadSlugOrID)
	if customErr != nil {
		return nil, customErr
	}

	if threadData.Title != "" {
		thread.Title = threadData.Title
	}
	if threadData.Message != "" {
		thread.Message = threadData.Message
	}

	if err := tu.threadRepo.Update(thread); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return thread, nil
}

func (tu *ThreadUsecase) GetBySlug(threadSlug string) (*models.Thread, *errors.Error) {
	thread, err := tu.threadRepo.SelectBySlug(threadSlug)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeThreadDoesNotExist, "slug", threadSlug)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return thread, nil
}

func (tu *ThreadUsecase) GetByID(threadID uint64) (*models.Thread, *errors.Error) {
	thread, err := tu.threadRepo.SelectByID(threadID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeThreadDoesNotExist, "id", strconv.Itoa(int(threadID)))
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return thread, nil
}

func (tu *ThreadUsecase) GetBySlugOrID(threadSlugOrID string) (*models.Thread, *errors.Error) {
	var thread *models.Thread
	var err *errors.Error

	threadID, parseErr := strconv.ParseUint(threadSlugOrID, 10, 64)
	threadSlug := threadSlugOrID
	if parseErr == nil {
		thread, err = tu.GetByID(threadID)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = tu.GetBySlug(threadSlug)
		if err != nil {
			return nil, err
		}
	}
	return thread, nil
}

func (tu *ThreadUsecase) GetByPostID(postID uint64) (*models.Thread, *errors.Error) {
	thread, err := tu.threadRepo.SelectByPostID(postID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeThreadDoesNotExist, "post id", strconv.Itoa(int(postID)))
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return thread, nil
}

func (tu *ThreadUsecase) CheckThreadExistence(threadSlugOrID string) (uint64, *errors.Error) {
	var threadID uint64
	var err error

	threadID, err = strconv.ParseUint(threadSlugOrID, 10, 64)
	threadSlug := threadSlugOrID
	if err != nil {
		threadID, err = tu.threadRepo.SelectIDBySlug(threadSlug)
		if err == sql.ErrNoRows {
			return 0, errors.BuildByMsg(CodeThreadDoesNotExist, "slug", threadSlug)
		}
	} else {
		_, err = tu.threadRepo.SelectIDByID(threadID)
		if err == sql.ErrNoRows {
			return 0, errors.BuildByMsg(CodeThreadDoesNotExist, "id", strconv.Itoa(int(threadID)))
		}
	}
	if err != nil {
		return 0, errors.New(CodeInternalError, err)
	}
	return threadID, nil
}

func (tu *ThreadUsecase) Vote(threadSlugOrID string, vote *models.Vote) (*models.Thread, *errors.Error) {
	threadID, customErr := tu.CheckThreadExistence(threadSlugOrID)
	if customErr != nil {
		return nil, customErr
	}

	if err := tu.threadRepo.VoteByID(threadID, vote); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}

	thread, customErr := tu.GetByID(threadID)
	if customErr != nil {
		return nil, customErr
	}
	return thread, nil
}

func (tu *ThreadUsecase) ListByForum(forumSlug string, since time.Time, pgnt *models.Pagination) ([]*models.Thread, *errors.Error) {
	threads, err := tu.threadRepo.SelectAllByForum(forumSlug, since, pgnt)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(threads) == 0 {
		return []*models.Thread{}, nil
	}
	return threads, nil
}
