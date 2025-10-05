package controllers

import (
	"errors"
	"net/mail"

	"github.com/claudesky/identity-go/utils"
)

// -- Structs

type LoginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type RegisterRequest struct {
	LoginRequest
}

type VerificationRequest struct {
	Id    *string `json:"id"`
	Token *string `json:"token"`
}

type RefreshRequest struct {
	RefreshToken *string `json:"refresh_token"`
}

type ResendVerificationRequest struct {
	Email *string `json:"email"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Message struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type DataMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

type ErrorMessage struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// -- Validators

func (rq *RegisterRequest) validate() (err error) {
	if err = rq.LoginRequest.validate(); err != nil {
		return
	}
	if _, err = mail.ParseAddress(*rq.Email); err != nil {
		return errors.New("[email] must be a valid email address")
	}
	if err = utils.ValidatePassword(*rq.Password); err != nil {
		return err
	}
	return
}

func (rq *LoginRequest) validate() error {
	if rq.Email == nil {
		return errors.New("[email] is required")
	}
	if rq.Password == nil {
		return errors.New("[password] is required")
	}
	return nil
}

func (rq VerificationRequest) validate() error {
	if rq.Id == nil {
		return errors.New("[id] is required")
	}
	if rq.Token == nil {
		return errors.New("[token] is required")
	}
	return nil
}

func (rq *RefreshRequest) validate() error {
	if rq.RefreshToken == nil {
		return errors.New("[refresh_token] is required")
	}
	return nil
}

func (rq *ResendVerificationRequest) validate() error {
	if rq.Email == nil {
		return errors.New("[email] is required")
	}
	if _, err := mail.ParseAddress(*rq.Email); err != nil {
		return errors.New("[email] must be a valid email address")
	}
	return nil
}
