package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// BanExists checks if a ban exists in the bans table for the specified users.
func (db *appdbimpl) BanExists(banningUsername, bannedUsername string) (bool, error) {
	// Query to check if the ban exists
	query := "SELECT EXISTS(SELECT 1 FROM Bans WHERE banningUser = ? AND bannedUser = ?)"

	// Execute the query
	var exists bool
	err := db.c.QueryRow(query, banningUsername, bannedUsername).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if ban exists: %w", err)
	}

	return exists, nil
}
