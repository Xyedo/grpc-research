package user

import (
	"context"
	"grpc-research/user-server/internal/entites"
)

type Usecase interface {
	AddUser(ctx context.Context, newUser entites.User) (string, error)
	GetUserById(ctx context.Context, id string) (entites.User, error)
	GetUser(ctx context.Context, userIdentifier entites.UserIdentifier) (entites.User, error)
}
