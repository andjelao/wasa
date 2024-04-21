package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) CommentExists(commentID int) (bool, error) {
	// Execute the query to check if the comment exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE comment_id = ?)", commentID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return false and an error indicating that the comment doesn't exist
			return false, err
		}
		// Return false and the error if there's an issue executing the query
		return false, fmt.Errorf("error checking if comment exists: %w", err)
	}
	// Return whether the comment exists or not
	return exists, nil
}
