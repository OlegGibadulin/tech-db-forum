package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
)

type UserPgRepository struct {
	dbConn *sql.DB
}

func NewUserPgRepository(conn *sql.DB) user.UserRepository {
	return &UserPgRepository{
		dbConn: conn,
	}
}

func (ur *UserPgRepository) Insert(user *models.User) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO users(nickname, fullname, email, about)
		VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.Email, user.About)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserPgRepository) Update(user *models.User) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = ur.dbConn.Exec(
		`UPDATE users
		SET fullname = $2, email = $3, about = $4
		WHERE nickname = $1;`,
		user.Nickname, user.Fullname, user.Email, user.About)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserPgRepository) SelectByNickname(nickname string) (*models.User, error) {
	user := &models.User{}

	row := ur.dbConn.QueryRow(
		`SELECT nickname, fullname, email, about
		FROM users
		WHERE nickname=$1`,
		nickname)

	err := row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserPgRepository) SelectByEmail(email string) (*models.User, error) {
	user := &models.User{}

	row := ur.dbConn.QueryRow(
		`SELECT nickname, fullname, email, about
		FROM users
		WHERE email=$1`,
		email)

	err := row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserPgRepository) SelectAllByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	rows, err := ur.dbConn.Query(
		`SELECT nickname, fullname, email, about
		FROM users
		WHERE nickname=$1 or email=$2`,
		nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserPgRepository) SelectAllByForum(forumSlug string, since string, pgnt *models.Pagination) ([]*models.User, error) {
	var values []interface{}

	selectQuery := `
		SELECT u.nickname, u.fullname, u.email, u.about
		FROM forum_user AS fu
		JOIN users AS u ON u.nickname=fu.nickname
		WHERE forum=$1`
	values = append(values, forumSlug)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY u.nickname DESC"
	} else {
		sortQuery = "ORDER BY u.nickname"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var filterQuery string
	if since != "" {
		ind := len(values) + 1
		filterQuery = "AND u.nickname > $" + strconv.Itoa(ind)
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		filterQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	rows, err := ur.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
