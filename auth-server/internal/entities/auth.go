package entities

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EmailLogin struct {
	Email    string
	Password string
}

var username = regexp.MustCompile(`^\S+$`)

func (login EmailLogin) Validate() error {
	return validation.ValidateStruct(&login,
		validation.Field(&login.Email, validation.Required, is.Email),
		validation.Field(&login.Password, validation.Required, validation.Length(5, 100)))
}

type UsernameLogin struct {
	Username string
	Password string
}

func (login UsernameLogin) Validate() error {
	return validation.ValidateStruct(&login,
		validation.Field(&login.Username, validation.Required, validation.Length(5, 20), validation.Match(username)),
		validation.Field(&login.Password, validation.Required, validation.Length(5, 100)))
}

type RefreshToken struct {
	RefreshToken string
}

func (logout RefreshToken) Validate() error {
	return validation.ValidateStruct(&logout, validation.Field(
		&logout.RefreshToken, validation.Required, validation.Length(10, 0)))
}
