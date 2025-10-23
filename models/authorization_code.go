package models

import (
	"time"

	"github.com/claudesky/identity-go/utils"
)

type AuthorizationCode struct {
	Id                  string    `json:"id"`
	Code                string    `json:"code"`
	ClientId            string    `json:"client_id"`
	UserId              string    `json:"user_id"`
	RedirectUri         string    `json:"redirect_uri"`
	Scope               string    `json:"scope"`
	ExpiresAt           time.Time `json:"expires_at"`
	CreatedAt           time.Time `json:"created_at"`
	Used                bool      `json:"used"`
	CodeChallenge       string    `json:"code_challenge,omitempty"`
	CodeChallengeMethod string    `json:"code_challenge_method,omitempty"`
}

func NewAuthorizationCode(
	clientId string,
	userId string,
	redirectUri string,
	scope string,
	codeChallenge string,
	codeChallengeMethod string,
) *AuthorizationCode {
	now := time.Now()
	return &AuthorizationCode{
		Id:                  utils.PseudoUUID(),
		Code:                utils.PseudoUUID(),
		ClientId:            clientId,
		UserId:              userId,
		RedirectUri:         redirectUri,
		Scope:               scope,
		ExpiresAt:           now.Add(time.Minute), // Authorization codes expire in 1 minute
		CreatedAt:           now,
		Used:                false,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}
