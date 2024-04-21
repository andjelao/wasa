package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) AddLike(like Like) error {
	// Execute the INSERT statement to add the like to the likes table
	_, err := db.c.Exec("INSERT INTO likes (username, photo_id) VALUES (?, ?)", like.Username, like.PhotoId)
	if err != nil {
		// Return the error if the INSERT operation fails
		return fmt.Errorf("error adding like to database: %w", err)
	}
	return nil
}
