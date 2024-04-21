package api

import (
	"net/http"
	"strings"
	// "fantastic-coffee-decaffeinated/service/api/reqcontext"
)

func GetUsernameFromToken(r *http.Request) string {
	// Get the Authorization header value
	authHeader := r.Header.Get("Authorization")

	// Check if the Authorization header is present and starts with "Bearer"
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		// Extract the token from the Authorization header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		// Assuming token directly represents the username
		return token
	}

	// Return an empty string if the Authorization header is missing or invalid
	return ""
}
