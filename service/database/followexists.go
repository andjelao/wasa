package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) FollowExists(follower, following string) (bool, error) {
	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM follow WHERE follower = ? AND followed = ?)"

	// Execute the query and scan the result
	var exists bool
	err := db.c.QueryRow(query, follower, following).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking follow existence: %w", err)
	}

	return exists, nil
}
