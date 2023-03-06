package usecase

import (
	"context"
	"grpc-research/auth-server/forum/user"
	auth "grpc-research/auth-server/internal"
	"grpc-research/auth-server/internal/entities"

	"golang.org/x/crypto/bcrypt"
)

func NewAuthUseCase(userGrpc user.UserClient, token auth.Tokenizer, productRepo auth.Repository) auth.UseCase {
	return &authUC{
		userClient: userGrpc,
		jwtToken:   token,
		authRepo:   productRepo,
	}
}

type authUC struct {
	userClient user.UserClient
	authRepo   auth.Repository
	jwtToken   auth.Tokenizer
}

// Logout implements auth.UseCase
func (uc *authUC) Logout(ctx context.Context, payload entities.RefreshToken) error {
	err := payload.Validate()
	if err != nil {
		return err
	}
	username, err := uc.jwtToken.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		return err
	}
	_, err = uc.userClient.GetUser(ctx, &user.GetUserByNameOrEmail{Unique: &user.GetUserByNameOrEmail_Username{Username: username}})
	if err != nil {
		return err
	}
	return nil
}

// RefreshAcess implements auth.UseCase
func (uc *authUC) RefreshAccess(ctx context.Context, payload entities.RefreshToken) (string, error) {
	err := payload.Validate()
	if err != nil {
		return "", err
	}
	username, err := uc.jwtToken.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		return "", err
	}
	accessToken, err := uc.jwtToken.GenerateAccessToken(username)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// LoginWithEmail implements auth.UseCase
func (uc *authUC) LoginWithEmail(ctx context.Context, payload entities.EmailLogin) (accessToken string, refreshToken string, err error) {
	err = payload.Validate()
	if err != nil {
		return "", "", err
	}
	userResp, err := uc.userClient.GetUser(ctx, &user.GetUserByNameOrEmail{Unique: &user.GetUserByNameOrEmail_Email{Email: payload.Email}})
	if err != nil {
		return "", "", err
	}
	return uc.login(userResp.GetUsername(), userResp.GetHashedpassword(), payload.Password)

}

// LoginWithUsername implements auth.UseCase
func (uc *authUC) LoginWithUsername(ctx context.Context, payload entities.UsernameLogin) (accessToken string, refreshToken string, err error) {
	err = payload.Validate()
	if err != nil {
		return "", "", err
	}
	userResp, err := uc.userClient.GetUser(ctx, &user.GetUserByNameOrEmail{Unique: &user.GetUserByNameOrEmail_Username{Username: payload.Username}})

	if err != nil {
		return "", "", err
	}
	return uc.login(userResp.GetUsername(), userResp.GetHashedpassword(), payload.Password)
}

func (uc *authUC) login(username, hashedPassword, password string) (accessToken string, refreshToken string, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", "", err
	}
	accessToken, err = uc.jwtToken.GenerateAccessToken(username)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = uc.jwtToken.GenerateRefreshToken(username)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
