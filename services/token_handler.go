package services

import (
	"crypto"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenHandler struct {
	pkey   crypto.PrivateKey
	pubkey crypto.PublicKey
}

func NewTokenHandler(
	pkey crypto.PrivateKey,
	pubkey crypto.PublicKey,
) *TokenHandler {
	return &TokenHandler{pkey: pkey, pubkey: pubkey}
}

func (th *TokenHandler) SignToken(claims jwt.Claims) (string, error) {
	// Sign a token and return the JWT
	tokenString, err := jwt.
		NewWithClaims(jwt.SigningMethodEdDSA, claims).
		SignedString(th.pkey)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func (th *TokenHandler) VerifyToken(ts string) (token *jwt.Token, err error) {

	token, err = jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return th.pubkey, nil
	})

	return
}
