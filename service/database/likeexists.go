package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) LikeExists(like Like) (bool, error) {
	// Execute the SELECT statement to check if the like exists in the likes table
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE username = ? AND photo_id = ?)", like.Username, like.PhotoId).Scan(&exists)
	if err != nil {
		// Return the error if the SELECT operation fails
		return false, fmt.Errorf("error checking if like exists in database: %w", err)
	}
	return exists, nil
}
