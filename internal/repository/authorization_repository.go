package repository

import (
	"context"
	"fmt"

	"github.com/arynskiii/help_desk/models"
	"github.com/jmoiron/sqlx"
)

type AuthMySQL struct {
	db *sqlx.DB
}

var empty models.User

func NewAuthMySQL(db sqlx.DB) *AuthMySQL {
	return &AuthMySQL{
		db: &db,
	}
}

func (auth *AuthMySQL) CreateUser(ctx context.Context, user models.User) (int64, error) {
	query := "INSERT INTO users (name,password) values (?,?)"
	res, err := auth.db.ExecContext(ctx, query, user.Name, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (auth *AuthMySQL) GetUser(name string) (models.User, error) {
	var user models.User
	query := "SELECT id, name, password FROM users WHERE name = ?"
	row := auth.db.QueryRowContext(context.Background(), query, name)
	err := row.Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return empty, fmt.Errorf("cannot scan user: %w", err)
	}
	return user, nil
}

func (a *AuthMySQL) SaveTokens(name string, token string) error {
	query := "UPDATE users SET token=? WHERE name=?"
	_, err := a.db.Exec(query, token, name)
	if err != nil {
		return fmt.Errorf("ERROR:don't save user's token: %w", err)
	}
	return nil
}

func (a *AuthMySQL) GetUserByToken(token string) (models.User, error) {
	var user models.User
	query := "select id,name,password from users where token=?"
	row := a.db.QueryRowContext(context.Background(), query, token)
	if err := row.Scan(&user.Id, &user.Name, &user.Password); err != nil {
		return empty, err
	}
	return user, nil
}
