//go:build !release

package main

import (
	"log/slog"

	"github.com/joho/godotenv"
)

// See build flag at top of file;
// This is only included when the `release` flag is not set on compile
func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Could not load .env file")
	}
}
