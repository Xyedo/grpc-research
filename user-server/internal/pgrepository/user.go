package pgrepository

import (
	"context"
	user "grpc-research/user-server/internal"
	"grpc-research/user-server/internal/entites"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewRepository(db *sqlx.DB) user.Repository {
	return &pgrepository{
		conn: db,
	}
}

type pgrepository struct {
	conn *sqlx.DB
}

func (db *pgrepository) GetUserById(ctx context.Context, id string) (entites.User, error) {
	var user entites.User
	err := db.conn.GetContext(ctx, &user, `SELECT id, username, email, password FROM users WHERE id = $1`, id)
	if err != nil {
		return entites.User{}, err
	}
	return user, nil
}

// GetUserByEmail implements user.Repository
func (db *pgrepository) GetUserByEmail(ctx context.Context, email string) (entites.User, error) {
	var user entites.User
	err := db.conn.GetContext(ctx, &user, `SELECT id, username, email, password FROM users WHERE email = $1`, email)
	if err != nil {
		return entites.User{}, err
	}
	return user, nil
}

// GetUserByUsername implements user.Repository
func (db *pgrepository) GetUserByUsername(ctx context.Context, username string) (entites.User, error) {
	var user entites.User
	err := db.conn.GetContext(ctx, &user, `SELECT id, username, email, password FROM users WHERE username = $1`, username)
	if err != nil {
		return entites.User{}, err
	}
	return user, nil
}

// AddUser implements user.Repository
func (db *pgrepository) AddUser(ctx context.Context, user entites.User) (string, error) {
	err := db.conn.GetContext(ctx, &user.Id, `INSERT INTO users (username, email, password) VALUES($1,$2,$3) RETURNING id `, strings.ToLower(user.Username), strings.ToLower(user.Email), user.Password)
	if err != nil {
		return "", err
	}
	return user.Id, nil

}
