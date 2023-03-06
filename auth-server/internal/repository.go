package auth

import (
	"context"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, refreshToken string) error
	IsRefreshTokenExists(ctx context.Context, refreshToken string) error
}
