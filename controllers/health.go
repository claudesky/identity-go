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

func (c *HealthController) Check(w http.ResponseWriter, _ *http.Request) {
	if c.healthy {
		respondWithMessage(w, "OK", http.StatusOK)
	} else {
		respondWithMessage(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}
