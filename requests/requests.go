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

type Authorize struct {
	ResponseType        *string `json:"response_type"`
	ClientId            *string `json:"client_id"`
	RedirectUri         *string `json:"redirect_uri"`
	Scope               *string `json:"scope"`
	State               *string `json:"state"`
	CodeChallenge       *string `json:"code_challenge"`
	CodeChallengeMethod *string `json:"code_challenge_method"`
}

type TokenExchange struct {
	GrantType    *string `json:"grant_type"`
	Code         *string `json:"code"`
	RedirectUri  *string `json:"redirect_uri"`
	ClientId     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
	CodeVerifier *string `json:"code_verifier"`
}

type CreateClient struct {
	Name         *string   `json:"name"`
	RedirectUris *[]string `json:"redirect_uris"`
}

type UpdateClient struct {
	Name         *string   `json:"name"`
	RedirectUris *[]string `json:"redirect_uris"`
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

func (rq *Authorize) Validate() error {
	if rq.ResponseType == nil {
		return errors.New("[response_type] is required")
	}
	if *rq.ResponseType != "code" {
		return errors.New("[response_type] must be 'code'")
	}
	if rq.ClientId == nil {
		return errors.New("[client_id] is required")
	}
	if rq.RedirectUri == nil {
		return errors.New("[redirect_uri] is required")
	}
	// PKCE: if code_challenge is provided, code_challenge_method must also be provided
	if rq.CodeChallenge != nil && rq.CodeChallengeMethod == nil {
		return errors.New("[code_challenge_method] is required when [code_challenge] is provided")
	}
	return nil
}

func (rq *TokenExchange) Validate() error {
	if rq.GrantType == nil {
		return errors.New("[grant_type] is required")
	}
	if *rq.GrantType != "authorization_code" {
		return errors.New("[grant_type] must be 'authorization_code'")
	}
	if rq.Code == nil {
		return errors.New("[code] is required")
	}
	if rq.RedirectUri == nil {
		return errors.New("[redirect_uri] is required")
	}
	if rq.ClientId == nil {
		return errors.New("[client_id] is required")
	}
	if rq.ClientSecret == nil {
		return errors.New("[client_secret] is required")
	}
	return nil
}

func (rq *CreateClient) Validate() error {
	if rq.Name == nil || *rq.Name == "" {
		return errors.New("[name] is required")
	}
	if rq.RedirectUris == nil || len(*rq.RedirectUris) == 0 {
		return errors.New("[redirect_uris] must contain at least one URI")
	}
	return nil
}

func (rq *UpdateClient) Validate() error {
	if rq.Name == nil || *rq.Name == "" {
		return errors.New("[name] is required")
	}
	if rq.RedirectUris == nil || len(*rq.RedirectUris) == 0 {
		return errors.New("[redirect_uris] must contain at least one URI")
	}
	return nil
}
