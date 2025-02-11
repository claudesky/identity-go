package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.
		GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidatePassword(password string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	if len(password) < 8 {
		return errors.New("[password] must be at least 8 characters long")
	}

	// Check each character type
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Ensure all conditions are met
	if !hasUpper {
		return errors.New("[password] must include at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("[password] must include at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("[password] must include at least one digit")
	}
	if !hasSpecial {
		return errors.New("[password] must include at least one special character")
	}

	return nil
}
