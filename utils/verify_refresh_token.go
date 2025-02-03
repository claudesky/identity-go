package utils

import (
	"log/slog"

	"github.com/claudesky/identity-go/services"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TYP string
	JTI string
	SUB string
	JTF string
}

func VerifyRefreshToken(th *services.TokenHandler, tokenString string) (bool, *Claims) {
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

	return true, &Claims{
		typ,
		jti,
		sub,
		jtf,
	}
}
