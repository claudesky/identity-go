package models

import "time"

type TokenFamily struct {
	Id           string    `json:"id"`
	Sub          string    `json:"sub"`
	LastIssued   string    `json:"last_issued"`
	CreatedAt    time.Time `json:"created_at"`
	LastIssuedAt time.Time `json:"last_issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
}
