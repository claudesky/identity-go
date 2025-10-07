package services

import (
	"crypto"
	"fmt"
	"log/slog"
	"time"

	"github.com/claudesky/identity-go/utils"
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
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.
				Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return th.pubkey, nil
	})

	return
}

func (th *TokenHandler) GenerateTokens(
	now *time.Time,
	sub string,
	jtf string,
	jtiRT string,
	additionalClaims jwt.MapClaims,
) (
	refreshString string,
	tokenString string,
	expRT time.Time,
	err error,
) {
	// Tokens TTL
	ttlRT := time.Hour * time.Duration(72)
	ttlAT := time.Minute * time.Duration(5)
	expRT = now.Add(ttlRT)
	expAT := now.Add(ttlAT)

	refreshString, err = th.SignToken(jwt.MapClaims{
		"jti": jtiRT,
		"jtf": jtf,
		"sub": sub,
		"exp": expRT.Unix(),
		"typ": "refresh_token",
	})
	if err != nil {
		slog.Info("refresh token signing failed", "error", err)
		return
	}

	// Build access token claims
	accessClaims := jwt.MapClaims{
		"jti": utils.PseudoUUID(),
		"jtf": jtf,
		"jtp": jtf,
		"sub": sub,
		"exp": expAT.Unix(),
		"typ": "access_token",
	}
	// Merge additional claims
	for k, v := range additionalClaims {
		accessClaims[k] = v
	}

	tokenString, err = th.SignToken(accessClaims)
	if err != nil {
		slog.Info("access token signing failed", "error", err)
		return
	}

	return
}

type verifyRefreshTokenOutput struct {
	TYP string
	JTI string
	SUB string
	JTF string
}

func (th *TokenHandler) VerifyRefreshToken(
	tokenString string,
) (
	bool,
	*verifyRefreshTokenOutput,
) {
	// Verify Token
	token, err := th.VerifyToken(tokenString)
	if err != nil {
		slog.Info("token verification failed", "error", err, "token", tokenString)
		return false, nil
	}

	// Parse claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("failed to parse token claims", "token", tokenString)
		return false, nil
	}

	// Check if correct refresh token type
	typ, ok := claims["typ"].(string)
	if !ok || typ != "refresh_token" {
		return false, nil
	}

	// Check claims are complete
	jti, ok := claims["jti"].(string)
	if !ok {
		slog.Error("no jti", "token", tokenString)
		return false, nil
	}

	sub, err := claims.GetSubject()
	if err != nil {
		slog.Error("no sub", "token", tokenString, "error", err)
		return false, nil
	}

	jtf, ok := claims["jtf"].(string)
	if !ok {
		slog.Error("no jtf", "token", tokenString)
		return false, nil
	}

	return true, &verifyRefreshTokenOutput{
		typ,
		jti,
		sub,
		jtf,
	}
}
