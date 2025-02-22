package models

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/claudesky/identity-go/utils"
)

type EmailVerificationRequest struct {
	Id                string    `json:"id"`
	RegisterRequestId *string   `json:"-"`
	Email             *string   `json:"email"`
	Token             *string   `json:"-"`
	Accepted          bool      `json:"-"`
	Revoked           bool      `json:"-"`
	ExpiresAt         time.Time `json:"expires_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// Default 5 min expiry
var evrExpiry = time.Minute * 5

func NewEmailVerificationRequest(
	rrid string,
	email string,
) *EmailVerificationRequest {
	code := strconv.Itoa(rand.Intn(900000) + 100000)
	now := time.Now()

	return &EmailVerificationRequest{
		Id:                utils.PseudoUUID(),
		RegisterRequestId: &rrid,
		Email:             &email,
		Token:             &code,
		Accepted:          false,
		Revoked:           false,
		ExpiresAt:         now.Add(evrExpiry),
		CreatedAt:         now,
	}
}
