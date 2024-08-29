package models

import "time"

type User struct {
	Id                    string     `json:"id"`
	Password              *string    `json:"-"`
	Name                  *string    `json:"name"`
	Email                 *string    `json:"email"`
	EmailVerifiedOn       *time.Time `json:"email_verified_on"`
	PhoneNumber           *string    `json:"phone_number"`
	PhoneNumberVerifiedOn *time.Time `json:"phone_number_verified_on"`
}
