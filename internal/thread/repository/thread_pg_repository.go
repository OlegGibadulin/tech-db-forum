package repository

import (
	"context"
	"database/sql"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
)

type ThreadPgRepository struct {
	dbConn *sql.DB
}

func NewThreadPgRepository(conn *sql.DB) thread.ThreadRepository {
	return &ThreadPgRepository{
		dbConn: conn,
	}
}

func (tr *ThreadPgRepository) Insert(thread *models.Thread) error {
	tx, err := tr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO threads(title, author, message, created, forum, slug)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, votes`,
		thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum, thread.Slug)

	err = row.Scan(&thread.ID, &thread.Votes)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tr *ThreadPgRepository) SelectBySlug(slug string) (*models.Thread, error) {
	thread := &models.Thread{}

	row := tr.dbConn.QueryRow(
		`SELECT id, title, author, message, created, forum, votes, slug
		FROM threads
		WHERE slug=$1`,
		slug)

	err := row.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Message, &thread.Created,
		&thread.Forum, &thread.Votes, &thread.Slug)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadPgRepository) SelectByID(threadID uint64) (*models.Thread, error) {
	thread := &models.Thread{}

	row := tr.dbConn.QueryRow(
		`SELECT id, title, author, message, created, forum, votes, slug
		FROM threads
		WHERE id=$1`,
		threadID)

	err := row.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Message, &thread.Created,
		&thread.Forum, &thread.Votes, &thread.Slug)
	if err != nil {
		return nil, err
	}
	return thread, nil
}
