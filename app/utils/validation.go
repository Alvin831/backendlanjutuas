package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername checks if username is valid (alphanumeric, 3-20 chars)
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword checks if password meets requirements (min 6 chars)
func ValidatePassword(password string) bool {
	return len(password) >= 6
}

// IsEmptyString checks if string is empty or only whitespace
func IsEmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}