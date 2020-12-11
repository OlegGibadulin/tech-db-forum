package repository

import (
	"context"
	"database/sql"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/service"
)

type ServicePgRepository struct {
	dbConn *sql.DB
}

func NewServicePgRepository(conn *sql.DB) service.ServiceRepository {
	return &ServicePgRepository{
		dbConn: conn,
	}
}

func (sr *ServicePgRepository) ClearAllTables() error {
	tx, err := sr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`TRUNCATE users, forums, forum_user, threads, posts, votes CASCADE`)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *ServicePgRepository) GetRowsCount() (*models.Status, error) {
	status := &models.Status{}

	row := sr.dbConn.QueryRow(
		`SELECT
		(SELECT COUNT(*) FROM users) as users_count,
		(SELECT COUNT(*) FROM forums) as forums_count,
		(SELECT COUNT(*) FROM threads) as threads_count,
		(SELECT COUNT(*) FROM posts) as posts_count`)

	err := row.Scan(&status.Users, &status.Forums, &status.Threads, &status.Posts)
	if err != nil {
		return nil, err
	}
	return status, nil
}
