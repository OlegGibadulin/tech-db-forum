package service

import (
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ServiceUsecase interface {
	Clear() *errors.Error
	GetStatus() (*models.Status, *errors.Error)
}
