package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/post"
)

type PostPgRepository struct {
	dbConn *sql.DB
}

func NewPostPgRepository(conn *sql.DB) post.PostRepository {
	return &PostPgRepository{
		dbConn: conn,
	}
}

func buildValuesQuery(rowsCount, valuesCount int) string {
	/*
		VALUES ($1, $2, ...),
		...
	*/
	var valuesQueries []string
	for i := 0; i < rowsCount; i++ {
		var values []string
		for j := 0; j < valuesCount; j++ {
			value := i*valuesCount + j + 1
			values = append(values, "$"+strconv.Itoa(value))
		}
		valuesQuery := fmt.Sprintf("(%s)", strings.Join(values, ", "))
		valuesQueries = append(valuesQueries, valuesQuery)
	}
	joinedValues := strings.Join(valuesQueries, ", ")
	return fmt.Sprintf("VALUES %s", joinedValues)
}

func (pr *PostPgRepository) Insert(posts []*models.Post, thread *models.Thread) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	var values []interface{}
	created := time.Now()

	selectQuery := "INSERT INTO posts(parent, author, message, forum, thread, created)"
	returnQuery := "RETURNING id, isedited, forum, thread, created"

	for _, post := range posts {
		values = append(values, post.Parent, post.Author, post.Message, thread.Forum, thread.ID, created)
	}
	valuesQuery := buildValuesQuery(len(posts), 6)

	resultQuery := strings.Join([]string{
		selectQuery,
		valuesQuery,
		returnQuery,
	}, " ")

	rows, err := tx.Query(resultQuery, values...)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	ind := 0
	for rows.Next() {
		err := rows.Scan(&posts[ind].ID, &posts[ind].IsEdited, &posts[ind].Forum,
			&posts[ind].Thread, &posts[ind].Created)
		if err != nil {
			tx.Rollback()
			return err
		}
		ind += 1
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (pr *PostPgRepository) Update(post *models.Post) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE posts
		SET message = $2
		WHERE id = $1`,
		post.ID, post.Message)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (pr *PostPgRepository) SelectByID(postID uint64) (*models.Post, error) {
	post := &models.Post{}

	row := pr.dbConn.QueryRow(
		`SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE id=$1`,
		postID)

	err := row.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
		&post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *PostPgRepository) SelectAllByThreadFlat(
	threadID uint64,
	since uint64,
	pgnt *models.Pagination) ([]*models.Post, error) {

	var values []interface{}

	selectQuery := `
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE thread=$1`
	values = append(values, threadID)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY id DESC"
	} else {
		sortQuery = "ORDER BY id"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var filterQuery string
	if since != 0 {
		ind := len(values) + 1
		if pgnt.Desc {
			filterQuery = "AND id < $" + strconv.Itoa(ind)
		} else {
			filterQuery = "AND id > $" + strconv.Itoa(ind)
		}
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		filterQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	rows, err := pr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostPgRepository) SelectAllByThreadTree(
	threadID uint64,
	since uint64,
	pgnt *models.Pagination) ([]*models.Post, error) {

	var values []interface{}

	selectQuery := `
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE thread=$1`
	values = append(values, threadID)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY path DESC"
	} else {
		sortQuery = "ORDER BY path"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var filterQuery string
	if since != 0 {
		ind := len(values) + 1
		subSelectQuery := fmt.Sprintf("(SELECT path FROM posts WHERE id=$%d)", ind)

		var subFilterQuery string
		if pgnt.Desc {
			subFilterQuery = "AND path <"
		} else {
			subFilterQuery = "AND path >"
		}

		filterQuery = strings.Join([]string{
			subFilterQuery,
			subSelectQuery,
		}, " ")
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		filterQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	rows, err := pr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func getSelectParentsQuery(
	threadID uint64,
	since uint64,
	pgnt *models.Pagination) (string, []interface{}) {

	var values []interface{}

	selectQuery := `
		SELECT id
		FROM posts
		WHERE thread=$1
		AND parent=0`
	values = append(values, threadID)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY id DESC"
	} else {
		sortQuery = "ORDER BY id"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var subSelectQuery string
	if since != 0 {
		subSelectQuery = `
		SELECT path[1]
		FROM posts
		WHERE id=$3`

		var filterQuery string
		if pgnt.Desc {
			filterQuery = "AND path[1] <"
		} else {
			filterQuery = "AND path[1] >"
		}
		subSelectQuery = fmt.Sprintf("%s (%s)", filterQuery, subSelectQuery)
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		subSelectQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	return resultQuery, values
}

func (pr *PostPgRepository) SelectAllByThreadParentTree(
	threadID uint64,
	since uint64,
	pgnt *models.Pagination) ([]*models.Post, error) {

	subSelectQuery, values := getSelectParentsQuery(threadID, since, pgnt)

	selectQuery := `
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE path[1] IN`

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY path[1] DESC, path, id"
	} else {
		sortQuery = "ORDER BY path, id"
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		fmt.Sprintf("(%s)", subSelectQuery),
		sortQuery,
	}, " ")

	rows, err := pr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
