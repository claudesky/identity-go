package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/utils"
	"github.com/jackc/pgx/v5"
)

type UserController struct {
	ur repositories.UserRepository
}

func NewUserController(
	ur repositories.UserRepository,
) *UserController {
	return &UserController{ur}
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.ur.GetAllUsers(r.Context())
	if err != nil {
		slog.Warn("failed to get all users", "error", err)
		internalServerError(w)
		return
	}

	respondWithData(w, "Success", users, 200)
}

func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithMessage(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := c.ur.GetUserById(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, utils.ErrInvalidUUID) {
		respondWithMessage(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Warn("failed to get user by id", "error", err, "id", id)
		internalServerError(w)
		return
	}

	json.NewEncoder(w).Encode(user)
}
