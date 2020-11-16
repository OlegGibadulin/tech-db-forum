package repository

import (
	"context"
	"database/sql"

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
