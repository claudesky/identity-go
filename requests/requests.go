package requests

import (
	"errors"
	"net/mail"

	"github.com/claudesky/identity-go/utils"
)

// -- Structs
type Login struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type Register struct {
	Login
}

type Verification struct {
	Id    *string `json:"id"`
	Token *string `json:"token"`
}

type Refresh struct {
	RefreshToken *string `json:"refresh_token"`
}

type ResendVerification struct {
	Email *string `json:"email"`
}

// -- Validators

func (rq *Login) Validate() error {
	if rq.Email == nil {
		return errors.New("[email] is required")
	}
	if rq.Password == nil {
		return errors.New("[password] is required")
	}
	return nil
}

func (rq *Register) Validate() (err error) {
	if err = rq.Login.Validate(); err != nil {
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

func (rq Verification) Validate() error {
	if rq.Id == nil {
		return errors.New("[id] is required")
	}
	if rq.Token == nil {
		return errors.New("[token] is required")
	}
	return nil
}

func (rq *Refresh) Validate() error {
	if rq.RefreshToken == nil {
		return errors.New("[refresh_token] is required")
	}
	return nil
}

func (rq *ResendVerification) Validate() error {
	if rq.Email == nil {
		return errors.New("[email] is required")
	}
	if _, err := mail.ParseAddress(*rq.Email); err != nil {
		return errors.New("[email] must be a valid email address")
	}
	return nil
}
