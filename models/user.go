package models

import (
	"time"

	"github.com/claudesky/identity-go/utils"
)

type User struct {
	Id                    string     `json:"id"`
	Password              *string    `json:"-"`
	Name                  *string    `json:"name"`
	Email                 *string    `json:"email"`
	EmailVerifiedOn       *time.Time `json:"email_verified_on"`
	PhoneNumber           *string    `json:"phone_number"`
	PhoneNumberVerifiedOn *time.Time `json:"phone_number_verified_on"`
}

func NewUser(password string, email string) *User {
	return &User{
		Id:       utils.PseudoUUID(),
		Password: &password,
		Email:    &email,
	}
}
