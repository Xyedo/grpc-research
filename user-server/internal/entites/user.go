package entites

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var username = regexp.MustCompile(`^\S+$`)

type UserIdentifier struct {
	Username string
	Email    string
}

func (u UserIdentifier) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.When(u.Username == "", validation.Required).Else(validation.Empty)),
		validation.Field(&u.Username, validation.When(u.Email == "", validation.Required).Else(validation.Empty)),
	)
}

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Username, validation.Required, validation.Match(username)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	)
}
