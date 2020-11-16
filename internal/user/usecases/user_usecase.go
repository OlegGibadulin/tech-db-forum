package usecases

import (
	"database/sql"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uu *UserUsecase) Create(user *models.User) *errors.Error {
	users, customErr := uu.ListByNicknameOrEmail(user.Nickname, user.Email)
	switch {
	case customErr != nil:
		return customErr
	case len(users) != 0:
		customErr = errors.BuildByBody(CodeUserAlreadyExists, users)
		return customErr
	}

	if err := uu.userRepo.Insert(user); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (uu *UserUsecase) GetByNickname(nickname string) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectByNickname(nickname)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeUserDoesNotExist, nickname)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return user, nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectByEmail(email)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.BuildByMsg(CodeUserDoesNotExist, email)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return user, nil
}

func (uu *UserUsecase) ListByNicknameOrEmail(nickname string, email string) ([]*models.User, *errors.Error) {
	users, err := uu.userRepo.SelectAllByNicknameOrEmail(nickname, email)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(users) == 0 {
		return []*models.User{}, nil
	}
	return users, nil
}
