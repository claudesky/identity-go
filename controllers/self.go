package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/services"
)

type SelfController struct {
	th  *services.TokenHandler
	ur  repositories.UserRepository
	tfr repositories.TokenFamilyRepository
}

func NewSelfController(
	th *services.TokenHandler,
	ur repositories.UserRepository,
	tfr repositories.TokenFamilyRepository,
) *SelfController {
	return &SelfController{th, ur, tfr}
}

func (c *SelfController) Sessions(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.th.VerifyToken(tokenString)
	if err != nil {
		slog.Info("token verification failed", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		slog.Warn("failed to get subject claim", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessions, err := c.tfr.GetTokensBySub(r.Context(), sub)
	if err != nil {
		slog.Warn("failed to get token families by subject", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessions)
}
