package usecases

import (
	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/service"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewServiceUsecase(repo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: repo,
	}
}

func (su *ServiceUsecase) Clear() *errors.Error {
	if err := su.serviceRepo.ClearAllTables(); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (su *ServiceUsecase) GetStatus() (*models.Status, *errors.Error) {
	status, err := su.serviceRepo.GetRowsCount()
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return status, nil
}
