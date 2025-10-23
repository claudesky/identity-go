package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/requests"
	"github.com/claudesky/identity-go/utils"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ClientController struct {
	cr repositories.ClientRepository
}

func NewClientController(
	cr repositories.ClientRepository,
) *ClientController {
	return &ClientController{cr}
}

// List all OAuth clients
func (c *ClientController) Index(w http.ResponseWriter, r *http.Request) {
	clients, err := c.cr.GetAllClients(r.Context())
	if err != nil {
		slog.Warn("failed to get all clients", "error", err)
		internalServerError(w)
		return
	}

	respondWithData(w, "Success", clients, http.StatusOK)
}

// Show a single OAuth client
func (c *ClientController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithMessage(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	client, err := c.cr.GetClientById(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, utils.ErrInvalidUUID) {
		respondWithMessage(w, "Client not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Warn("failed to get client by id", "error", err, "id", id)
		internalServerError(w)
		return
	}

	respondWithData(w, "Success", client, http.StatusOK)
}

// Create a new OAuth client
func (c *ClientController) Create(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.CreateClient,
) {
	// Generate client secret
	clientSecret := utils.PseudoUUID()

	// Hash the client secret
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash client secret", "error", err)
		internalServerError(w)
		return
	}

	// Create client model
	now := time.Now()
	client := models.NewClient(*rq.Name, *rq.RedirectUris, string(hashedSecret))
	client.CreatedAt = now
	client.UpdatedAt = now

	// Save to database
	err = c.cr.InsertClient(r.Context(), client)
	if err != nil {
		slog.Error("failed to insert client", "error", err)
		internalServerError(w)
		return
	}

	// Return client with plaintext secret (only time it will be shown)
	respondWithData(w, "Client created successfully", map[string]interface{}{
		"client":        client,
		"client_secret": clientSecret,
	}, http.StatusCreated)
}

// Update an existing OAuth client
func (c *ClientController) Update(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.UpdateClient,
) {
	id := r.PathValue("id")
	if id == "" {
		respondWithMessage(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Get existing client
	client, err := c.cr.GetClientById(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, utils.ErrInvalidUUID) {
		respondWithMessage(w, "Client not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Warn("failed to get client by id", "error", err, "id", id)
		internalServerError(w)
		return
	}

	// Update client fields
	client.Name = *rq.Name
	client.RedirectUris = *rq.RedirectUris
	client.UpdatedAt = time.Now()

	// Save to database
	err = c.cr.UpdateClient(r.Context(), client)
	if err != nil {
		slog.Error("failed to update client", "error", err, "id", id)
		internalServerError(w)
		return
	}

	respondWithData(w, "Client updated successfully", client, http.StatusOK)
}

// Delete an OAuth client
func (c *ClientController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithMessage(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Verify client exists
	_, err := c.cr.GetClientById(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, utils.ErrInvalidUUID) {
		respondWithMessage(w, "Client not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Warn("failed to get client by id", "error", err, "id", id)
		internalServerError(w)
		return
	}

	// Delete client
	err = c.cr.DeleteClient(r.Context(), id)
	if err != nil {
		slog.Error("failed to delete client", "error", err, "id", id)
		internalServerError(w)
		return
	}

	respondWithMessage(w, "Client deleted successfully", http.StatusOK)
}
