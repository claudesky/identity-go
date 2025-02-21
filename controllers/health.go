package controllers

import (
	"net/http"
)

type HealthController struct {
	healthy bool
}

func NewHealthController() *HealthController {
	return &HealthController{healthy: true}
}

func (c *HealthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health/check", c.check)
}

func (c *HealthController) check(w http.ResponseWriter, _ *http.Request) {
	if c.healthy {
		respondWithMessage(w, "OK", http.StatusOK)
	} else {
		respondWithError(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}
