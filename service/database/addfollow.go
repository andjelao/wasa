package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) AddFollow(follower, following string) error {
	// Prepare the SQL statement
	query := "INSERT INTO follow (follower, followed) VALUES (?, ?)"

	// Execute the SQL statement
	_, err := db.c.Exec(query, follower, following)
	if err != nil {
		return fmt.Errorf("error adding follow: %w", err)
	}

	return nil
}
