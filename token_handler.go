package main

import (
	"crypto"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenHandler struct {
	pkey   crypto.PrivateKey
	pubkey crypto.PublicKey
}

func newTokenHandler(pkey crypto.PrivateKey, pubkey crypto.PublicKey) *TokenHandler {
	return &TokenHandler{
		pkey:   pkey,
		pubkey: pubkey,
	}
}

func fetchPkey() (crypto.PrivateKey, error) {
	// Get the Private Key
	pkeyBytes, err := os.ReadFile(*privKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pkey, err := jwt.ParseEdPrivateKeyFromPEM(pkeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pkey, nil
}

func fetchPubkey() (crypto.PublicKey, error) {
	// Get the Public Key
	pubkeyBytes, err := os.ReadFile(*pubKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pubkey, err := jwt.ParseEdPublicKeyFromPEM(pubkeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pubkey, nil
}

func signToken(th *TokenHandler, claims jwt.Claims) (string, error) {
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

func verifyToken(th *TokenHandler, ts string) (*jwt.Token, error) {

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return th.pubkey, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}
