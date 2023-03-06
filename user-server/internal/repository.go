package user

import (
	"context"
	"grpc-research/user-server/internal/entites"
)

type Repository interface {
	AddUser(ctx context.Context, user entites.User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (entites.User, error)
	GetUserByEmail(ctx context.Context, email string) (entites.User, error)
	GetUserById(ctx context.Context, id string) (entites.User, error)
}
