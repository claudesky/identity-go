package models

import (
	"math/rand"
	"strconv"
	"time"
)

type EmailVerificationRequest struct {
	RegisterRequestId *string   `json:"register_request_id"`
	Email             *string   `json:"email"`
	Token             *string   `json:"token"`
	Accepted          bool      `json:"accepted"`
	Revoked           bool      `json:"revoked"`
	ExpiresAt         time.Time `json:"expires_at"`
	CreatedAt         time.Time `json:"created_at"`
}

var evrExpiry = time.Minute * 5

func NewEmailVerificationRequest(
	rrid string,
	email string,
) *EmailVerificationRequest {
	code := strconv.Itoa(rand.Intn(900000) + 100000)
	now := time.Now().UTC()

	return &EmailVerificationRequest{
		RegisterRequestId: &rrid,
		Email:             &email,
		Token:             &code,
		Accepted:          false,
		Revoked:           false,
		ExpiresAt:         now.Add(evrExpiry),
		CreatedAt:         now,
	}
}
