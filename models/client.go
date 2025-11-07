package models

import (
	"time"

	"github.com/claudesky/identity-go/utils"
)

type Client struct {
	Id           string    `json:"id"`
	ClientSecret string    `json:"-"`
	Name         string    `json:"name"`
	RedirectUris []string  `json:"redirect_uris"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewClient(name string, redirectUris []string, clientSecret string) *Client {
	now := time.Now()
	return &Client{
		Id:           utils.PseudoUUID(),
		ClientSecret: clientSecret,
		Name:         name,
		RedirectUris: redirectUris,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
