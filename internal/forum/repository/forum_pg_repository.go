package repository

import (
	"context"
	"database/sql"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
)

type ForumPgRepository struct {
	dbConn *sql.DB
}

func NewForumPgRepository(conn *sql.DB) forum.ForumRepository {
	return &ForumPgRepository{
		dbConn: conn,
	}
}

func (fr *ForumPgRepository) Insert(forum *models.Forum) error {
	tx, err := fr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO forums(title, author, slug)
		VALUES ($1, $2, $3)
		RETURNING posts, threads`,
		forum.Title, forum.User, forum.Slug)

	err = row.Scan(&forum.Posts, &forum.Threads)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (fr *ForumPgRepository) SelectBySlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}

	row := fr.dbConn.QueryRow(
		`SELECT title, author, slug, posts, threads
		FROM forums
		WHERE slug=$1`,
		slug)

	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (fr *ForumPgRepository) SelectByPostID(postID uint64) (*models.Forum, error) {
	forum := &models.Forum{}

	row := fr.dbConn.QueryRow(
		`SELECT title, author, slug, posts, threads
		FROM forums AS f
		JOIN posts AS p ON p.forum=f.slug
		WHERE p.id=$1`,
		postID)

	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, err
	}
	return forum, nil
}
