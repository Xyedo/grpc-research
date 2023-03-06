package usecase

import (
	"context"
	user "grpc-research/user-server/internal"
	"grpc-research/user-server/internal/entites"

	"golang.org/x/crypto/bcrypt"
)

func NewUserUseCase(userRepo user.Repository) user.Usecase {
	return &userUC{
		userRepo: userRepo,
	}
}

type userUC struct {
	userRepo user.Repository
}

func (uc *userUC) GetUser(ctx context.Context, userIdentifier entites.UserIdentifier) (entites.User, error) {
	err := userIdentifier.Validate()
	if err != nil {
		return entites.User{}, err
	}
	if userIdentifier.Email != "" {
		userDb, err := uc.userRepo.GetUserByEmail(ctx, userIdentifier.Email)
		if err != nil {
			return entites.User{}, err
		}
		return userDb, err
	}
	userDb, err := uc.userRepo.GetUserByUsername(ctx, userIdentifier.Username)
	if err != nil {
		return entites.User{}, err
	}
	return userDb, nil
}

func (uc *userUC) GetUserById(ctx context.Context, id string) (entites.User, error) {
	userDb, err := uc.userRepo.GetUserById(ctx, id)
	if err != nil {
		return entites.User{}, err
	}
	return userDb, nil
}

func (uc *userUC) AddUser(ctx context.Context, newUser entites.User) (string, error) {
	err := newUser.Validate()
	if err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	newUser.Password = string(hashed)
	id, err := uc.userRepo.AddUser(ctx, newUser)
	if err != nil {
		return "", err
	}
	return id, nil
}
