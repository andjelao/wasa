package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) AddComment(comment Comment) (int64, error) {
	// novo unused
	// Execute the INSERT statement to add the new comment to the comments table
	result, err := db.c.Exec("INSERT INTO comments ( photo_id, username, text, date) VALUES (?, ?, ?, ?)",
		comment.PhotoId, comment.Author, comment.Text, comment.Date)
	if err != nil {
		// Return the error if the INSERT operation fails
		return 0, fmt.Errorf("error adding comment to database: %w", err)
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		// Return the error if the INSERT operation fails
		return 0, fmt.Errorf("error adding comment to database: %w", err)
	}
	return lastInsertedID, nil
}
