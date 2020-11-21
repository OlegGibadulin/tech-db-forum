package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

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

func (tr *ThreadPgRepository) Update(thread *models.Thread) error {
	tx, err := tr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tr.dbConn.Exec(
		`UPDATE threads
		SET title = $2, message = $3
		WHERE id = $1;`,
		thread.ID, thread.Title, thread.Message)
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

func (tr *ThreadPgRepository) SelectByPostID(postID uint64) (*models.Thread, error) {
	thread := &models.Thread{}

	row := tr.dbConn.QueryRow(
		`SELECT id, title, author, message, created, forum, votes, slug
		FROM threads AS t
		JOIN posts AS p ON p.thread=t.id
		WHERE p.id=$1`,
		postID)

	err := row.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Message, &thread.Created,
		&thread.Forum, &thread.Votes, &thread.Slug)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadPgRepository) SelectAllByForum(forumSlug string, since time.Time, pgnt *models.Pagination) ([]*models.Thread, error) {
	var values []interface{}

	selectQuery := `
		SELECT id, title, author, message, created, forum, votes, slug
		FROM threads
		WHERE forum=$1`
	values = append(values, forumSlug)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY created DESC"
	} else {
		sortQuery = "ORDER BY created"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var filterQuery string
	if !since.IsZero() {
		ind := len(values) + 1
		filterQuery = "AND created >= $" + strconv.Itoa(ind)
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		filterQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	rows, err := tr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*models.Thread
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Message, &thread.Created,
			&thread.Forum, &thread.Votes, &thread.Slug)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return threads, nil
}
