package models

import "time"

type RegisterRequest struct {
	Id        string    `json:"id"`
	Password  *string   `json:"-"`
	Email     *string   `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
