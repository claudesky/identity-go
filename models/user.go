package models

import "time"

type User struct {
	Id                    string
	Password              *string `json:"-"`
	Name                  *string
	Email                 *string
	EmailVerifiedOn       *time.Time
	PhoneNumber           *string
	PhoneNumberVerifiedOn *time.Time
}
