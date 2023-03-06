package auth

import (
	"context"
	"grpc-research/auth-server/internal/entities"
)

type UseCase interface {
	LoginWithUsername(ctx context.Context, payload entities.UsernameLogin) (accessToken, refreshToken string, err error)
	LoginWithEmail(ctx context.Context, payload entities.EmailLogin) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, payload entities.RefreshToken) error
	RefreshAccess(ctx context.Context, payload entities.RefreshToken) (string, error)
}
