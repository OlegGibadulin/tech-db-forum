package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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
			value := rowsCount*valuesCount + j + 1
			values = append(values, strconv.Itoa(value))
		}
		valuesQuery := fmt.Sprintf("VALUES (%s)", strings.Join(values, ", "))
		valuesQueries = append(valuesQueries, valuesQuery)
	}
	return strings.Join(valuesQueries, ", ")
}

func (pr *PostPgRepository) Insert(posts []*models.Post, threadID uint64) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	var values []interface{}

	selectQuery := "INSERT INTO posts(parent, author, message, forum, thread)"
	returnQuery := "RETURNING id, isedited, created"

	for _, post := range posts {
		values = append(values, post.Parent, post.Author, post.Message, post.Forum, threadID)
	}
	valuesQuery := buildValuesQuery(len(posts), 5)

	resultQuery := strings.Join([]string{
		selectQuery,
		valuesQuery,
		returnQuery,
	}, " ")

	rows, err := tx.Query(resultQuery, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	ind := 0
	for rows.Next() {
		err := rows.Scan(&posts[ind].ID, &posts[ind].IsEdited, &posts[ind].Created)
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

	_, err = pr.dbConn.Exec(
		`UPDATE posts
		SET message = $2
		WHERE id = $1;`,
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

func (pr *PostPgRepository) SelectNotExistingParentPosts(posts []*models.Post) ([]uint64, error) {
	/*
		SELECT newp.parent FROM
		(VALUES (1), (2), (3)) as newp(parent)
		LEFT OUTER JOIN posts as oldp on oldp.id=newp.parent
		WHERE oldp.id IS NULL
	*/
	var values []interface{}

	var parentPosts []string
	for _, post := range posts {
		if post.Parent != 0 { // 0 is root
			values = append(values, post.Parent)
			ind := len(values)
			parentPost := fmt.Sprintf("($%d)", strconv.Itoa(ind))
			parentPosts = append(parentPosts, parentPost)
		}
	}
	valuesQuery := strings.Join(parentPosts, ", ") // ($1), ($2), ($3), ...
	parentPostsQuery := fmt.Sprintf("(VALUES %s) as newp(parent)", valuesQuery)

	selectQuery := "SELECT newp.parent FROM"
	diffQuery := `
		LEFT OUTER JOIN posts as oldp on oldp.id=newp.parent
		WHERE oldp.id IS NULL`

	resultQuery := strings.Join([]string{
		selectQuery,
		parentPostsQuery,
		diffQuery,
	}, " ")

	rows, err := pr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postsID []uint64
	for rows.Next() {
		var postID uint64
		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		postsID = append(postsID, postID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return postsID, nil
}
