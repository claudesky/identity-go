package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidUUID = errors.New("invalid UUID format")

// Match standard UUID format (8-4-4-4-12 hex digits)
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// Note - NOT RFC4122 compliant
func PseudoUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = strings.ToLower(
		fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
	)

	return
}

// Note - Slow?
func ValidateUUID(uuid string) error {
	if !uuidRegex.MatchString(strings.ToLower(uuid)) {
		return ErrInvalidUUID
	}
	return nil
}
