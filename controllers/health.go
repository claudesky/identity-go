package controllers

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

type HealthController struct {
	healthy bool
}

func NewHealthController() *HealthController {
	return &HealthController{healthy: true}
}

func (c *HealthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health/hello-word", c.helloWorld)
	mux.HandleFunc("GET /health/check", c.check)
}

func (c *HealthController) check(w http.ResponseWriter, _ *http.Request) {
	if c.healthy {
		json.NewEncoder(w).Encode(&Message{Message: "ok"})
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}

func (c *HealthController) helloWorld(w http.ResponseWriter, _ *http.Request) {
	if c.healthy {
		json.NewEncoder(w).Encode(&Message{Message: "Hello, world!"})
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}
