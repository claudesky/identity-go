package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/middleware"
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
	sub := middleware.GetSubject(r)

	sessions, err := c.tfr.GetTokensBySub(r.Context(), sub)
	if err != nil {
		slog.Warn("failed to get token families by subject", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessions)
}
