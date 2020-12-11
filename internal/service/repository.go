package service

import "github.com/OlegGibadulin/tech-db-forum/internal/models"

type ServiceRepository interface {
	ClearAllTables() error
	GetRowsCount() (*models.Status, error)
}
