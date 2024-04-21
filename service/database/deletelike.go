package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) DeleteLike(like Like) error {
	// Execute the DELETE query to remove the like from the database
	_, err := db.c.Exec("DELETE FROM likes WHERE photo_id = ? AND username = ?", like.PhotoId, like.Username)
	if err != nil {
		// Return the error if the deletion fails
		return fmt.Errorf("error deleting like from database: %w", err)
	}
	return nil
}
