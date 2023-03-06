package pgrepository

import (
	"context"
	auth "grpc-research/auth-server/internal"

	"github.com/jmoiron/sqlx"
)

func NewAuthRepo(conn *sqlx.DB) auth.Repository {
	return &authRepo{
		conn: conn,
	}
}

type authRepo struct {
	conn *sqlx.DB
}

// IsRefreshTokenExists implements auth.Repository
func (a *authRepo) IsRefreshTokenExists(ctx context.Context, refreshToken string) error {
	var retRefreshToken string
	err := a.conn.GetContext(ctx, &retRefreshToken, `SELECT token FROM authentications WHERE refresh_token = $1`, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

// CreateRefreshToken implements auth.Repository
func (a *authRepo) CreateRefreshToken(ctx context.Context, refreshToken string) error {
	var retRefreshToken string
	err := a.conn.GetContext(ctx, &retRefreshToken, `INSERT INTO authentications (token) VALUES ($1) RETURNING token`, refreshToken)
	if err != nil {
		return err
	}
	return nil
}
