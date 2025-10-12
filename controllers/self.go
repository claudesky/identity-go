package controllers

import (
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

func (c *SelfController) Self(w http.ResponseWriter, r *http.Request) {
	sub := middleware.GetSubject(r)

	user, err := c.ur.GetUserById(r.Context(), sub)
	if err != nil {
		slog.Warn("failed to get user by id", "error", err)
		internalServerError(w)
		return
	}

	respondWithData(w, "Success", user, 200)
}

func (c *SelfController) Sessions(w http.ResponseWriter, r *http.Request) {
	sub := middleware.GetSubject(r)

	sessions, err := c.tfr.GetTokensBySub(r.Context(), sub)
	if err != nil {
		slog.Warn("failed to get token families by subject", "error", err)
		internalServerError(w)
		return
	}

	respondWithData(w, "Success", sessions, 200)
}
