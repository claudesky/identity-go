package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func DecodeRequestJSON[T any](r *http.Request) (*T, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("could not decode request", "error", err)
		return nil, fmt.Errorf("invalid request payload")
	}

	return &req, nil
}
